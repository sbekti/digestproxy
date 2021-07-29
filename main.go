package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	digest "github.com/bobziuchkovski/digest"
	cli "github.com/urfave/cli"
)

var (
	port      int
	upstream  string
	proxyUser string
	proxyPass string
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:        "port",
			Value:       8080,
			Usage:       "port number to listen on",
			EnvVar:      "PORT",
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "upstream",
			Value:       "",
			Usage:       "upstream",
			EnvVar:      "UPSTREAM",
			Destination: &upstream,
		},
		cli.StringFlag{
			Name:        "username",
			Value:       "",
			Usage:       "upstream username",
			EnvVar:      "PROXY_USER",
			Destination: &proxyUser,
		},
		cli.StringFlag{
			Name:        "password",
			Value:       "",
			Usage:       "upstream password",
			EnvVar:      "PROXY_PASS",
			Destination: &proxyPass,
		},
	}

	app.Action = func(c *cli.Context) error {
		origin, _ := url.Parse(upstream)

		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)
			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}

		proxy := &httputil.ReverseProxy{
			Director:  director,
			Transport: digest.NewTransport(proxyUser, proxyPass),
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		})

		fmt.Printf("Proxy listening on port %d\n", port)
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
