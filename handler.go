package main

import (
	"flag"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	cacheControl = flag.String("cache-control", "private,no-cache,max-age=0", "The Cache-Control header to return for all responses.")
	contentType  = flag.String("content-type", "binary/octet-stream", "The Content-Type to return for all responses.")
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
	total := first + count*size

	w.Header().Set("Cache-Control", *cacheControl)
	w.Header().Set("Content-Length", strconv.Itoa(total))
	w.Header().Set("Content-Type", *contentType)
	w.WriteHeader(http.StatusOK)

	if writeRand(w, first) != nil {
		return
	}

	var ch <-chan time.Time
	if sleep > 0 {
		t := time.NewTicker(sleep)
		defer t.Stop() // Avoid leaking the ticker.

		ch = t.C
	} else {
		// closed channels are always ready, which will behave like
		// a zero delay.
		cc := make(chan time.Time)
		close(cc)

		ch = cc
	}

	ctx := r.Context()
	for i := 0; i < count; i++ {
		select {
		case <-ch:
			if writeRand(w, size) != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// formInt returns the integer value of a http.request form parameter, or 0
// if the parmeter is missing or not a number.
func formInt(r *http.Request, param string) int {
	if v, _ := strconv.Atoi(r.FormValue(param)); v > 0 && v <= 1000000000 {
		return v
	}
	return 0
}

// writeRand writes n random bytes to w then flushes the result.
func writeRand(w http.ResponseWriter, n int) error {
	// Copy n bytes from rand.Read to w without allocating.
	if _, err := io.CopyN(w, readerFunc(rand.Read), int64(n)); err != nil {
		return err
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}

// readerFunc is a function type that implements io.Reader
type readerFunc func(p []byte) (n int, err error)

func (r readerFunc) Read(p []byte) (n int, err error) {
	return r(p)
}
