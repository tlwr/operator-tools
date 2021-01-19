package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/urfave/cli/v2"

	cl "github.com/tlwr/operator-tools/pkg/colour"
	"github.com/tlwr/operator-tools/pkg/timeline"
)

func HTTPCmd() *cli.Command {
	return &cli.Command{
		Name: "http",
		Subcommands: []*cli.Command{
			{
				Name:  "profile",
				Usage: "Profiles an HTTP request, printing timing and response headers",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "url",
						Aliases:  []string{"u"},
						Required: true,
					},
				},

				Action: func(c *cli.Context) error {
					rawURL := c.String("url")
					if rawURL == "" {
						cli.ShowAppHelpAndExit(c, 1)
					}

					parsedURL, err := url.ParseRequestURI(c.String("url"))
					if err != nil {
						return err
					}

					trace := HTTPTrace{}

					ctx := httptrace.WithClientTrace(context.TODO(), trace.Trace())

					req, err := http.NewRequestWithContext(
						ctx,
						"GET",
						parsedURL.String(),
						/* body */ nil,
					)
					if err != nil {
						return err
					}

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						fmt.Println(err)
						return err
					}

					if trace.TLSConnStateInfo != nil {
						cs := trace.TLSConnStateInfo

						fmt.Println("TLS:")

						fmt.Printf("%s: ", cl.Blue("Version"))
						if cs.Version == tls.VersionSSL30 || cs.Version == tls.VersionTLS10 || cs.Version == tls.VersionTLS11 {
							versions := map[uint16]string{
								tls.VersionSSL30: "SSLv3",
								tls.VersionTLS10: "TLSv1.0",
								tls.VersionTLS11: "TLSv1.1",
							}
							fmt.Println(cl.Red(versions[cs.Version]))
						} else if cs.Version == tls.VersionTLS12 {
							fmt.Println(cl.Blue("TLSv1.2"))
						} else if cs.Version == tls.VersionTLS13 {
							fmt.Println(cl.Green("TLSv1.3"))
						} else {
							fmt.Println(cl.Yellow(fmt.Sprintf("%v", cs.Version)))
						}

						fmt.Printf("%s: ", cl.Blue("Cipher-Suite"))
						suiteIsInsecure := false
						suiteIsFound := false
						for _, suite := range tls.InsecureCipherSuites() {
							if suite.ID == cs.CipherSuite {
								suiteIsInsecure = true
								break
							}
						}
						for _, suite := range tls.CipherSuites() {
							if suite.ID == cs.CipherSuite {
								suiteIsFound = true
								break
							}
						}
						if suiteIsInsecure {
							fmt.Println(cl.Red(tls.CipherSuiteName(cs.CipherSuite)))
						} else if suiteIsFound {
							fmt.Println(cl.Green(tls.CipherSuiteName(cs.CipherSuite)))
						} else {
							fmt.Println(cl.Yellow(fmt.Sprintf("%v", cs.CipherSuite)))
						}

						fmt.Printf("%s: %s\n", cl.Blue("Server-Name"), cl.Yellow(cs.ServerName))
						fmt.Printf("%s: %s\n", cl.Blue("Negotiated-Protocol"), cl.Yellow(cs.NegotiatedProtocol))

						fmt.Println()
					}

					fmt.Println("Status:")
					fmt.Printf(
						"%s %s\n",
						cl.Blue(fmt.Sprintf(
							"HTTP/%d.%d",
							resp.ProtoMajor,
							resp.ProtoMinor,
						)),
						cl.Yellow(resp.Status),
					)

					fmt.Println("\nHeaders:")
					for headerName, headerVals := range resp.Header {
						for _, headerVal := range headerVals {
							fmt.Printf(
								"%s: %s\n",
								cl.Blue(headerName),
								cl.Yellow(headerVal),
							)
						}
					}

					_, err = ioutil.ReadAll(resp.Body)
					if err != nil {
						return err
					}

					trace.Finish()

					width, _, err := terminal.GetSize(0)
					if err != nil {
						return err
					}
					width /= 2

					entries := []timeline.TimelineEntry{}
					entries = append(
						entries,
						timeline.NewTimelineEntry(trace.DNSLookupStart, &trace.DNSLookupDone, "dns"),
						timeline.NewTimelineEntry(trace.ConnectStart, &trace.ConnectDone, "connect"),
					)
					if trace.TLSConnStateInfo != nil {
						entries = append(entries, timeline.NewTimelineEntry(trace.TLSStart, &trace.TLSDone, "tls"))
					}
					entries = append(
						entries,
						timeline.NewTimelineEntry(trace.WroteRequestHeadersDone, nil, "request-headers-done"),
						timeline.NewTimelineEntry(trace.WroteRequestDone, nil, "request-done"),
						timeline.NewTimelineEntry(trace.FirstResponseByteDone, &trace.End, "reading-response"),
					)

					tl := timeline.NewTimeline(trace.Start, trace.End, width, entries...)

					rendered, err := tl.Render()
					if err != nil {
						return err
					}

					fmt.Println("\nTrace:")
					fmt.Print(rendered)

					return nil
				},
			},
		},
	}
}

type HTTPTrace struct {
	Start time.Time
	End   time.Time

	DNSLookupStart          time.Time
	DNSLookupDone           time.Time
	ConnectStart            time.Time
	ConnectDone             time.Time
	TLSStart                time.Time
	TLSDone                 time.Time
	WroteRequestHeadersDone time.Time
	WroteRequestDone        time.Time
	FirstResponseByteDone   time.Time

	TLSConnStateInfo *tls.ConnectionState
}

func (t *HTTPTrace) Trace() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(_ string) {
			t.Start = time.Now()
		},

		DNSStart: func(_ httptrace.DNSStartInfo) {
			t.DNSLookupStart = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			t.DNSLookupDone = time.Now()
		},

		ConnectStart: func(_, _ string) {
			t.ConnectStart = time.Now()
		},
		ConnectDone: func(_, _ string, _ error) {
			t.ConnectDone = time.Now()
		},

		TLSHandshakeStart: func() {
			t.TLSStart = time.Now()
		},
		TLSHandshakeDone: func(cs tls.ConnectionState, _ error) {
			t.TLSDone = time.Now()

			t.TLSConnStateInfo = &cs
		},

		WroteHeaders: func() {
			t.WroteRequestHeadersDone = time.Now()
		},

		WroteRequest: func(_ httptrace.WroteRequestInfo) {
			t.WroteRequestDone = time.Now()
		},

		GotFirstResponseByte: func() {
			t.FirstResponseByteDone = time.Now()
		},
	}
}

func (t *HTTPTrace) Finish() {
	t.End = time.Now()
}
