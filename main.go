package main

import (
	"fmt"
	"github.com/ndrewnee/reverse-proxy/proxy"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: reverse-proxy <host> <search> <replace>")
		return
	}

	host := os.Args[1]
	search := os.Args[2]
	replace := os.Args[3]

	reverseProxy, err := proxy.NewDummyProxy(host, search, replace)
	if err != nil {
		log.Fatal("Parse host error:", err)
	}

	port := ":3000"
	log.Println("Started server on", port)

	err = http.ListenAndServe(port, reverseProxy)
	if err != nil {
		log.Fatal("Listen server error: ", err)
	}
}
