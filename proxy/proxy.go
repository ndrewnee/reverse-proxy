package proxy

import (
	"fmt"
	"github.com/ndrewnee/reverse-proxy/transport"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(host, search, replace string) http.Handler {
	proxyUrl, err := url.Parse(fmt.Sprintf("https://%s", host))
	if err != nil {
		log.Fatalf("Parse url '%s' error: %s", host, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy.Transport = transport.NewTransport(search, replace)

	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = req.URL.Host
	}

	return proxy
}
