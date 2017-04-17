package main

import (
	"flag"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"strings"
)

var (
	enableHttps = flag.Bool("enable-https", false, "Enable listening on port 443 for HTTPS connections and fetching of LetsEncrypt certificates")
	hostnames   = flag.String("hostnames", "", "A comma-separated whitelist of domains to try asking LetsEncrypt for an TLS cert (unset = any)")
	listen      = flag.String("listen", ":80", "[IP]:port to listen for HTTP connections.")
)

func main() {
	flag.Parse()

	if *enableHttps {
		go serveHttps()
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func serveHttps() {
	https := http.NewServeMux()
	https.HandleFunc("/", handler)
	whitelist := []string{}
	if *hostnames != "" {
		whitelist = strings.Split(*hostnames, ",")
	}
	log.Fatal(http.Serve(autocert.NewListener(whitelist...), https))
}
