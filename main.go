package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

var (
	enableHTTPS = flag.Bool("enable-https", false, "Enable listening on port 443 for HTTPS connections and fetching of LetsEncrypt certificates")
	hostnames   = flag.String("hostnames", "", "A comma-separated allowlist of domains to try asking LetsEncrypt for an TLS cert (unset = any)")
	listen      = flag.String("listen", ":80", "[IP]:port to listen for HTTP connections.")
)

func main() {
	flag.Parse()

	m := http.NewServeMux()
	m.HandleFunc("/", handler)

	if *enableHTTPS {
		go serveHTTPS(m)
	}

	serveHTTP(m)
}

func serveHTTP(handler http.Handler) {
	l, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(serve(l, handler))
}

func serveHTTPS(handler http.Handler) {
	allowlist := []string{}
	if *hostnames != "" {
		allowlist = strings.Split(*hostnames, ",")
	}

	log.Fatal(serve(autocert.NewListener(allowlist...), handler))
}

func serve(l net.Listener, handler http.Handler) error {
	s := &http.Server{
		Handler:           handler,
		ReadTimeout:       10 * time.Minute,
		WriteTimeout:      10 * time.Minute,
		ReadHeaderTimeout: 20 * time.Second,
	}

	return s.Serve(l)
}
