package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	conn, err := net.Dial("tcp", ":1443")
	if err != nil {
		log.Fatal(err)
	}

	tr := http2.Transport{}
	c, err := tr.NewClientConn(conn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("CanTakeNewRequest", c.CanTakeNewRequest())
	if err := c.Ping(context.TODO()); err != nil {
		log.Fatal(err)
	}

	for {
		req, err := http.NewRequest("CONNECT", "http://google.com:443", nil)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := c.RoundTrip(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		v, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(v))

		time.Sleep(1000 * time.Millisecond)
	}
}
