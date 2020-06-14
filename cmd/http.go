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
					trace.Begin()

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

					fmt.Printf("%s %s\n", resp.Proto, resp.Status)
					fmt.Println("\nHeaders:")
					for headerName, headerVals := range resp.Header {
						for _, headerVal := range headerVals {
							fmt.Printf("%s: %s\n", headerName, headerVal)
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

					tl := timeline.NewTimeline(
						trace.Start, trace.End, width,

						timeline.NewTimelineEntry(
							trace.DNSLookupStart, &trace.DNSLookupDone, "dns",
						),

						timeline.NewTimelineEntry(
							trace.ConnectStart, &trace.ConnectDone, "connect",
						),

						timeline.NewTimelineEntry(
							trace.TLSStart, &trace.TLSDone, "tls",
						),

						timeline.NewTimelineEntry(
							trace.WroteRequestHeadersDone, nil, "request-headers-done",
						),

						timeline.NewTimelineEntry(
							trace.WroteRequestDone, nil, "request-done",
						),

						timeline.NewTimelineEntry(
							trace.FirstResponseByteDone, &trace.End, "reading-response",
						),
					)

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
}

func (t *HTTPTrace) Trace() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
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
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			t.TLSDone = time.Now()
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

func (t *HTTPTrace) Begin() {
	t.Start = time.Now()
}

func (t *HTTPTrace) Finish() {
	t.End = time.Now()
}
