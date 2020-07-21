package cmd

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	textutil "golang.org/x/tools/godoc/util"

	"github.com/urfave/cli/v2"
)

func X509Cmd() *cli.Command {
	var (
		expiryDays = 30
		excludeStr = ".git,node_modules,vendor,fixture,testdata"

		emptySpaceRegex = regexp.MustCompile(`\s*`)
		nakedCertRegex  = regexp.MustCompile(`(?m)MII[\sA-Za-z0-9+\/]*[=]*`)
	)

	return &cli.Command{
		Name: "x509",
		Subcommands: []*cli.Command{
			{
				Name:    "find-certs",
				Aliases: []string{"fc"},
				Usage:   "Finds x509 certificates recursively",

				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "expiry-days",
						Destination: &expiryDays,
						Value:       expiryDays,
					},

					&cli.StringFlag{
						Name:        "exclude",
						Destination: &excludeStr,
						Value:       excludeStr,
					},
				},

				Action: func(c *cli.Context) error {
					expiryDuration := time.Hour * time.Duration(24*expiryDays)

					excludedPathFrags := []string{}
					for _, s := range strings.Split(excludeStr, ",") {
						exclude := strings.TrimSpace(s)
						if len(exclude) > 0 {
							excludedPathFrags = append(excludedPathFrags, exclude)
						}
					}

					textFiles := []string{}
					err := filepath.Walk(
						".",
						func(path string, info os.FileInfo, err error) error {
							if info.IsDir() {
								return nil
							}

							if !info.Mode().IsRegular() {
								return nil
							}

							for _, excluded := range excludedPathFrags {
								if strings.Contains(path, excluded) {
									return nil
								}
							}

							contents, err := ioutil.ReadFile(path)
							if err != nil {
								fmt.Println(path, err)
								return err
							}

							if textutil.IsText(contents) {
								textFiles = append(textFiles, path)
							}

							return nil
						},
					)
					if err != nil {
						return err
					}

					textFilesWithCerts := []FileRawCertificate{}
					for _, filepath := range textFiles {
						contents, err := ioutil.ReadFile(filepath)
						if err != nil {
							return err
						}

						nakedCertMatches := nakedCertRegex.FindAll(contents, -1)
						for _, match := range nakedCertMatches {
							textFilesWithCerts = append(
								textFilesWithCerts,
								FileRawCertificate{
									Filepath:       filepath,
									RawCertificate: match,
								},
							)
						}
					}

					parsedCertificates := []FileCertificate{}
					for _, fileCert := range textFilesWithCerts {
						pem := emptySpaceRegex.ReplaceAll(fileCert.RawCertificate, []byte{})

						// FIXME byte -> str -> byte -> str gross but works
						der, err := base64.StdEncoding.DecodeString(string(pem))
						if err != nil {
							continue
						}

						certificate, err := x509.ParseCertificate([]byte(der))
						if err != nil {
							continue
						}

						parsedCertificates = append(parsedCertificates, FileCertificate{
							Filepath:    fileCert.Filepath,
							Certificate: certificate,
							RawPEM:      pem,
						})
					}

					expiryTime := time.Now().Add(expiryDuration)
					expCertificatesByFile := make(map[string][]FileCertificate)
					for _, parsedFileCert := range parsedCertificates {
						if parsedFileCert.Certificate.NotAfter.Before(expiryTime) {
							expCertificatesByFile[parsedFileCert.Filepath] = append(
								expCertificatesByFile[parsedFileCert.Filepath],
								parsedFileCert,
							)
						}
					}

					for file, expCerts := range expCertificatesByFile {
						fmt.Println(file)

						sort.Slice(expCerts, func(i int, j int) bool {
							iCert := expCerts[i].Certificate
							jCert := expCerts[j].Certificate
							return iCert.NotAfter.Before(jCert.NotAfter)
						})

						for _, cert := range expCerts {
							expiresInHours := cert.Certificate.NotAfter.Sub(time.Now()).Hours()

							if expiresInHours <= 0 {
								fmt.Printf(
									"- %s (%s) expired %d days ago\n",
									cert.Certificate.Issuer,
									string(cert.RawPEM[0:32]),
									-1*int(expiresInHours/24),
								)
							} else {
								fmt.Printf(
									"- %s (%s) expires in %d days\n",
									cert.Certificate.Issuer,
									string(cert.RawPEM[0:32]),
									int(expiresInHours/24),
								)
							}
						}
					}

					return nil
				},
			},
		},
	}
}

type FileRawCertificate struct {
	Filepath       string
	RawCertificate []byte
}

type FileCertificate struct {
	Filepath    string
	RawPEM      []byte
	Certificate *x509.Certificate
}
