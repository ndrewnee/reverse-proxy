package proxy

import (
	"fmt"
	"github.com/ndrewnee/reverse-proxy/transport"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(host, search, replace string) (http.Handler, error) {
	proxyUrl, err := url.Parse(fmt.Sprintf("https://www.%s", host))
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy.Transport = transport.NewTransport(search, replace)

	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = req.URL.Host
	}

	return proxy, nil
}
