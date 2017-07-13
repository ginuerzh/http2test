package main

import (
	"log"
	"net"
	"net/http"

	"golang.org/x/net/http2"
)

func init() {
	http2.VerboseLogs = true
}

func main() {
	server := http2.Server{}

	ln, err := net.Listen("tcp", ":1443")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("accept", conn.RemoteAddr())

		go func() {
			opt := http2.ServeConnOpts{
				Handler: http.HandlerFunc(h2cHandler),
			}
			server.ServeConn(conn, &opt)
		}()
	}
}

func h2cHandler(w http.ResponseWriter, r *http.Request) {
	r.Write(w)
}
