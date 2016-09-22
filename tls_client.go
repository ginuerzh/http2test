package main

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	tr := http2.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MaxVersion:         tls.VersionTLS12,
			MinVersion:         tls.VersionTLS12,
			CipherSuites:       []uint16{tls.VersionTLS12},
		},
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			log.Println(cfg.MaxVersion, cfg.MinVersion)
			return tls.DialWithDialer(&net.Dialer{Timeout: 30 * time.Second}, "tcp", "localhost:8080", cfg)
		},
	}
	client := http.Client{Transport: &tr}
	req := &http.Request{
		Method: "CONNECT",
		URL:    &url.URL{Scheme: "https"},
		Host:   "www.example.com:22",
		Header: make(http.Header),
	}
	req.Header.Set("Proxy-Connection", "keep-alive")

	dump, _ := httputil.DumpRequest(req, false)
	log.Println(string(dump))

	for {
		go func() {
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			dump, _ = httputil.DumpResponse(resp, false)
			log.Println(string(dump))
		}()

		time.Sleep(50 * time.Millisecond)
	}
}
