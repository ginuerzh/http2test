package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/http2"
)

func init() {
	http2.VerboseLogs = true
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	conn, err := tls.Dial("tcp", ":1443", &tls.Config{
		InsecureSkipVerify:       true,
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Handshake(); err != nil {
		log.Fatal()
	}

	tr := http2.Transport{}
	c, err := tr.NewClientConn(conn)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Ping(context.TODO()); err != nil {
		log.Fatal(err)
	}

	// stream(c)

	for {
		for i := 0; i < 100; i++ {
			go stream(c)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func stream(cc *http2.ClientConn) {
	pr, pw := io.Pipe()

	log.Println("CanTakeNewRequest", cc.CanTakeNewRequest())

	go func() {
		//for {
		//	time.Sleep(5 * time.Microsecond)
		fmt.Fprintf(pw, "It is now %v\n", time.Now())
		//}
		pw.Close()
	}()

	req, err := http.NewRequest("CONNECT", "http://google.com:443", ioutil.NopCloser(pr))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := cc.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	exit := make(chan interface{})
	go func() {
		n, err := io.Copy(os.Stdout, resp.Body)
		if err != nil {
			log.Fatalf("copied %d, %v", n, err)
		}
		close(exit)
	}()

	<-exit
	//log.Println("exited")
	// select {}
}
