package main

import (
	"flag"
	"log"
	"net"
	"net/http"
)

func main() {
	devMode := flag.Bool("dev", false, "Activate developer mode")
	flag.Parse()
	port := ":8080"
	handler := http.FileServer(http.Dir("dist"))

	if *devMode {
		http.Handle("/", devNoCacheMiddleware(loggingMiddleware(handler)))
	} else {
		http.Handle("/", loggingMiddleware(handler))
	}
	log.Printf("Server running: %s%s", getOutboundIp(), port)

	http.ListenAndServe(port, nil)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func( w http.ResponseWriter, r *http.Request) {
			log.Printf("[%s] %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
	})
}

func devNoCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
		    // TODO: should it not be set on w?
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			next.ServeHTTP(w, r)
	})
}

func getOutboundIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
