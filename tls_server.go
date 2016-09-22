package main

import (
	//"fmt"
	"crypto/tls"
	"golang.org/x/net/http2"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

func main() {
	http2.VerboseLogs = true
	cert, err := tls.LoadX509KeyPair("localhost.cert", "localhost.key")
	if err != nil {
		log.Fatal(err)
	}
	ln, err := tls.Listen("tcp", ":8080",
		&tls.Config{
			Certificates: []tls.Certificate{cert},
			MaxVersion:   tls.VersionTLS12,
			MinVersion:   tls.VersionTLS12,
			CipherSuites: []uint16{tls.VersionTLS12},
		})
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}
func ShowRequestInfoHandler(w http.ResponseWriter, r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	log.Println(string(dump))
}

func handleHttp2Proxy(w http.ResponseWriter, r *http.Request) {
	ShowRequestInfoHandler(w, r)
	time.Sleep(60 * time.Second)
	w.WriteHeader(http.StatusOK)
}

func handleConn(conn net.Conn) {
	server := http2.Server{PermitProhibitedCipherSuites: true}
	server.ServeConn(conn, &http2.ServeConnOpts{Handler: http.HandlerFunc(handleHttp2Proxy)})
}
