package main

import (
	"net"
	"net/http"
	"log"
)

func main() {
	handler := http.FileServer(http.Dir("dist"))
	http.Handle("/", loggingMiddleware(handler))

	log.Printf("Server running: %s", getOutboundIp())

	http.ListenAndServe(":8080", nil)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func( w http.ResponseWriter, r *http.Request) {
			log.Printf("[%s] %s", r.Method, r.URL.Path)
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
