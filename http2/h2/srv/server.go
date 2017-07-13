package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func init() {
	http2.VerboseLogs = true
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}

	config := tls.Config{
		Certificates:             []tls.Certificate{cert},
		PreferServerCipherSuites: true,
	}

	log.Fatal(http2Server(&config))
}

func httpServer(config *tls.Config) error {
	srv := http.Server{
		Addr:      ":1443",
		Handler:   http.HandlerFunc(http2Handler),
		TLSConfig: config,
	}
	http2.ConfigureServer(&srv, nil)
	return srv.ListenAndServeTLS("", "")
}

func http2Server(config *tls.Config) error {
	server := http2.Server{}

	ln, err := tls.Listen("tcp", ":1443", config)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		log.Println("accept", conn.RemoteAddr())

		if tc, ok := conn.(*tls.Conn); ok {
			if err := tc.Handshake(); err != nil {
				log.Fatal(err)
			}
		}

		go func() {
			opt := http2.ServeConnOpts{
				Handler: http.HandlerFunc(http2Handler),
			}
			server.ServeConn(conn, &opt)
		}()
	}
}

func http2Handler(w http.ResponseWriter, r *http.Request) {
	io.Copy(flushWriter{w}, r.Body)
}

type flushWriter struct {
	w io.Writer
}

func (fw flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if f, ok := fw.w.(http.Flusher); ok {
		f.Flush()
	}
	return
}
