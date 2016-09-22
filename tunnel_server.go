package main

import (
	//"fmt"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

func main() {
	var srv http.Server
	http2.VerboseLogs = true
	srv.Addr = ":8080"
	srv.Handler = http.HandlerFunc(handleHttp2Proxy)
	// This enables http2 support
	http2.ConfigureServer(&srv, nil)
	log.Fatal(srv.ListenAndServeTLS("localhost.cert", "localhost.key"))
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
