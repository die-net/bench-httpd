package main

import (
	"flag"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	cacheControl = flag.String("cache-control", "private,no-cache,max-age=0", "The Cache-Control header to return for all responses.")
	contentType  = flag.String("content-type", "text/plain", "The Content-Type to return for all responses.")
)

func handler(w http.ResponseWriter, r *http.Request) {
	// first: Count of bytes to send before any sleep
	first := formInt(r, "first")
	// sleep: Sleep in ms converted to time.Duration
	sleep := time.Duration(formInt(r, "sleep")) * time.Millisecond
	// size: Count of bytes to send after sleep
	size := formInt(r, "size")
	// count: Count of sleep+send cycles
	count := formInt(r, "count")
	if count <= 0 && (sleep > 0 || size > 0) {
		count = 1
	}

	w.Header().Set("Content-Type", *contentType)
	w.Header().Set("Cache-Control", *cacheControl)
	w.WriteHeader(http.StatusOK)

	if writeRand(w, first) != nil {
		return
	}

	for i := 0; i < count; i++ {
		time.Sleep(sleep)
		if writeRand(w, size) != nil {
			return
		}
	}
}

// formInt returns the integer value of a http.request form parameter, or 0
// if the parmeter is missing or not a number.
func formInt(r *http.Request, param string) int {
	v, _ := strconv.Atoi(r.FormValue(param))
	return v
}

// writeRand writes n random bytes to w then flushes the result.
func writeRand(w http.ResponseWriter, n int) error {
	if n <= 0 {
		return nil
	}

	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	if _, err := w.Write([]byte(b)); err != nil {
		return err
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}
