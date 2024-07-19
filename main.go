package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/acme/autocert"
)

var (
	enableHTTPS = flag.Bool("enable-https", false, "Enable listening on port 443 for HTTPS connections and fetching of LetsEncrypt certificates")
	hostnames   = flag.String("hostnames", "", "A comma-separated allowlist of domains to try asking LetsEncrypt for an TLS cert (unset = any)")
	listen      = flag.String("listen", ":80", "[IP]:port to listen for HTTP connections.")
)

func main() {
	flag.Parse()

	if *enableHTTPS {
		go serveHTTPS()
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func serveHTTPS() {
	https := http.NewServeMux()
	https.HandleFunc("/", handler)
	allowlist := []string{}
	if *hostnames != "" {
		allowlist = strings.Split(*hostnames, ",")
	}
	log.Fatal(http.Serve(autocert.NewListener(allowlist...), https))
}
