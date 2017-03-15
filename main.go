package main

import (
	"github.com/ndrewnee/reverse-proxy/proxy"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: reverse-proxy <host> <search> <replace>")
	}

	host := os.Args[1]
	search := os.Args[2]
	replace := os.Args[3]

	reverseProxy := proxy.NewReverseProxy(host, search, replace)

	port := ":3000"
	log.Println("Started server on", port)

	err := http.ListenAndServe(port, reverseProxy)
	if err != nil {
		log.Fatal("Listen server error: ", err)
	}
}
