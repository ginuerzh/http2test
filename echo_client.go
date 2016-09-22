package main

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	tr := http2.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return tls.DialWithDialer(&net.Dialer{Timeout: 30 * time.Second}, network, addr, cfg)
		},
	}
	client := http.Client{Transport: &tr}

	pr, pw := io.Pipe()
	req, err := http.NewRequest("PUT", "https://localhost:8080", ioutil.NopCloser(pr))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Fprintf(pw, "It is now %v\n", time.Now())
		}
	}()
	go func() {
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Got: %#v", res)
		n, err := io.Copy(os.Stdout, res.Body)
		log.Fatalf("copied %d, %v", n, err)
	}()
	select {}
}
