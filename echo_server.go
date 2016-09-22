package main

import (
	"golang.org/x/net/http2"
	"io"
	"log"
	"net/http"
)

func main() {
	var srv http.Server
	http2.VerboseLogs = true
	srv.Addr = ":8080"
	srv.Handler = http.HandlerFunc(echoCapitalHandler)
	http2.ConfigureServer(&srv, nil)
	log.Fatal(srv.ListenAndServeTLS("localhost.cert", "localhost.key"))
}

type capitalizeReader struct {
	r io.Reader
}

func (cr capitalizeReader) Read(p []byte) (n int, err error) {
	n, err = cr.r.Read(p)
	for i, b := range p[:n] {
		if b >= 'a' && b <= 'z' {
			p[i] = b - ('a' - 'A')
		}
	}
	return
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

func echoCapitalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "PUT required.", 400)
		return
	}
	io.Copy(flushWriter{w}, capitalizeReader{r.Body})
}
