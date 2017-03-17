// This utility start reverse proxy server that can search
// string in response body of host and replace it to the given string
// Usage: reverse-proxy <host> <search> <replace>
// <host> must contain "http" or "https"
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
		fmt.Printf("Usage: %s <host> <search> <replace>.\n<host> must contain 'http' or 'https'\n", os.Args[0])
		return
	}

	host := os.Args[1]
	search := os.Args[2]
	replace := os.Args[3]

	reverseProxy, err := proxy.NewReverseProxy(host, search, replace)
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
