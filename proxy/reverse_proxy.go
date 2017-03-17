// Package contains simple version of Reverse Proxy Server.
// It search string in response body and replace it to the given string
package proxy

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

var (
	ErrHostNotHttpOrHttps = errors.New(`host must start with "http" or "https"`)
)

// NewReverseProxy returns new reverse proxy handler that replaces text in response
func NewReverseProxy(host, search, replace string) (http.Handler, error) {
	if !strings.Contains(host, "http") {
		return nil, ErrHostNotHttpOrHttps
	}

	hostUrl, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(hostUrl)
	reverseProxy.Transport = &http.Transport{DialTLS: dialTLS, DisableCompression: true}

	director := reverseProxy.Director
	reverseProxy.Director = func(req *http.Request) {
		director(req)
		// Setting host because default director cannot do that
		req.Host = req.URL.Host
		// Don't send client's headers to real site
		req.Header = http.Header{}
	}

	reverseProxy.ModifyResponse = func(resp *http.Response) (err error) {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		err = resp.Body.Close()
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

// dialTLS is custom TLS dialer to verify host
func dialTLS(network, addr string) (net.Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{ServerName: host}

	tlsConn := tls.Client(conn, cfg)
	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		return nil, err
	}

	cs := tlsConn.ConnectionState()
	cert := cs.PeerCertificates[0]

	cert.VerifyHostname(host)

	return tlsConn, nil
}
