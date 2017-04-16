# bench-httpd

A tiny Go-based httpd server suitable for calibrating benchmarks and some networking testing.

When accessed over HTTP, it generates a response based on the supplied query parameters:

    first - Count of bytes to send before any sleep
    sleep - Sleep in milliseconds
    size  - Count of bytes to send after sleep
    count - Count of sleep+send cycles

All values are assumed to be 0 if not set. If either ```sleep``` or ```size``` are greater than 0, count is set to at least 1.

Which are used according to the following pseudo-code:

* Write ```first``` random bytes
* Loop ```count``` times, each loop doing:
  * Delay for ```sleep``` milliseconds
  * Write ```size``` random bytes

For example, accessing ```http://localhost/?first=1000&sleep=250&size=128&count=5``` will write ```1000``` random bytes, and loop ```5``` times that each: sleep ```250ms``` and write ```128``` bytes.

### Building

If you haven't used Go before (assuming you are on a Mac):

    brew install go
    mkdir -p $HOME/go/src
    export GOPATH=$HOME/go

Then to fetch the code and its dependencies and build the binary:

    go get -u -t github.com/die-net/bench-httpd

The binary is now available as: ```$HOME/go/bin/bench-httpd```

### Usage

    $HOME/go/bin/bench-httpd [flags]

Where ```[flags]``` are any of:

    -cache-control string
        The Cache-Control header to return for all responses. (default "private,no-cache,max-age=0")
    -content-type string
        The Content-Type to return for all responses. (default "text/plain")
    -enable-https
        Enabling listening on port 443 for HTTPS connections and fetching of Let's Encrypt certificates
    -hostnames string
        A comma-separated whitelist of hostnames to try asking Let's Encrypt for a TLS cert (unset = any)
    -listen string
        [IP]:port to listen for HTTP connections. (default ":80")

The ```-enable-https``` flag will attempt to get a new TLS certificate from Let's Encrypt for any hostname that it is accessed as.  This will only work for servers that are publicly reachable on port 443, because the Let's Encrypt validation will connect there before issuing the certificate.  If the ```-hostnames``` flag is set, only the listed hostnames are attempted, and TLS connections for all other hostnames are rejected.