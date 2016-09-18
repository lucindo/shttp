package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/facebookgo/httpdown"
	"github.com/fiorix/go-web/httpxtra"
	"github.com/pkg/browser"
	"github.com/rs/cors"
)

// HTTP Handler to disable caching

type noCacheLogHandler struct {
	handler http.Handler
}

func noCacheLogHandlerServer(handler http.Handler) http.Handler {
	return &noCacheLogHandler{handler}
}

func (s *noCacheLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	headers.Set("Pragma", "no-cache")
	headers.Set("Expires", "0")
	s.handler.ServeHTTP(w, r)
}

// HTTP Handler for Request log

type requestLogHandler struct {
	handler http.Handler
}

func requestLogServer(handler http.Handler) http.Handler {
	return &requestLogHandler{handler}
}

func (s *requestLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Proto: %v", r.Proto)
	s.handler.ServeHTTP(w, r)
}

// HTTP Handler for Apache log

func logApache(r *http.Request, created time.Time, status, bytes int) {
	fmt.Println(httpxtra.ApacheCommonLog(r, created, status, bytes))
}

func logHandler(handler http.Handler, logger httpxtra.LoggerFunc, xheader bool) http.Handler {
	return &httpxtra.Handler{Handler: handler, Logger: logger, XHeaders: xheader}
}

// Aux. functions

func formatAddress(host string, port int) (string, string) {
	address := fmt.Sprintf("%s:%d", host, port)
	return address, fmt.Sprintf("http://%s/", address)
}

func main() {
	// Server options
	port := flag.Int("port", 8080, "Port to bind the server")
	host := flag.String("host", "localhost", "Listen address")
	readTimeout := flag.Duration("rtimeout", 10*time.Second, "Server read timeout")
	writeTimeout := flag.Duration("wtimeout", 10*time.Second, "Server write timeout")
	maxHeaderBytes := flag.Int("maxheaders", 1<<20, "Max header size in bytes")

	// General Options
	dir := flag.String("dir", ".", "Directory to expose")
	open := flag.Bool("open", false, "Open a browser pointing to this server")
	quiet := flag.Bool("quiet", false, "Do not log requests")
	disableCors := flag.Bool("nocors", false, "Disable CORS headers")
	logRequest := flag.Bool("debug", false, "Log request information")
	noCache := flag.Bool("nocache", false, "Add headers disabling HTTP cache for all requests")

	flag.Parse()

	var handler http.Handler

	handler = http.FileServer(http.Dir(*dir))

	if !*quiet {
		handler = logHandler(handler, logApache, true)
	}
	if !*disableCors {
		handler = cors.Default().Handler(handler)
	}
	if *logRequest {
		if *quiet {
			fmt.Println("incompactible options 'debug' and 'quiet', exiting.")
			return
		}
		handler = requestLogServer(handler)
	}
	if *noCache {
		handler = noCacheLogHandlerServer(handler)
	}

	address, url := formatAddress(*host, *port)

	server := &http.Server{
		Addr:           address,
		Handler:        handler,
		ReadTimeout:    *readTimeout,
		WriteTimeout:   *writeTimeout,
		MaxHeaderBytes: *maxHeaderBytes,
	}

	// Provides graceful shutdown
	hd := &httpdown.HTTP{
		StopTimeout: 1 * time.Second,
		KillTimeout: 1 * time.Second,
	}

	if !*quiet {
		fmt.Printf("Listening to %s\n", url)
		fmt.Println("Hit CTRL-C to exit...")
	}

	// The browser may open fast enough (before we start ListenAndServe)
	// so I need to implement the httdown.ListenAndServe on a separete
	// gorroutine and deal with signals.
	if *open {
		err := browser.OpenURL(url)
		if err != nil {
			fmt.Printf("error opening browser: %v", err)
		}
	}

	if err := httpdown.ListenAndServe(server, hd); err != nil {
		panic(err)
	}
}
