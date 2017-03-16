package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func validateHost(host string) (string, error) {
	if !strings.Contains(host, "http") {
		host = "http://" + host
	}

	_, err := url.ParseRequestURI(host)
	return host, err
}

func NewReverseProxy(host, search, replace string) (http.Handler, error) {
	host, err := validateHost(host)
	if err != nil {
		return nil, err
	}

	hostUrl, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(hostUrl)

	director := reverseProxy.Director
	reverseProxy.Director = func(req *http.Request) {
		director(req)
		req.Host = req.URL.Host
	}

	reverseProxy.ModifyResponse = func(resp *http.Response) (err error) {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		b = bytes.Replace(b, []byte(search), []byte(replace), -1)
		body := ioutil.NopCloser(bytes.NewReader(b))

		resp.Body = body
		resp.ContentLength = int64(len(b))
		resp.Header.Set("Content-Length", strconv.Itoa(len(b)))

		return
	}

	return reverseProxy, nil
}
