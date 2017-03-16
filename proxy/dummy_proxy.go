package proxy

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func NewDummyProxy(host, search, replace string) (*DummyProxy, error) {
	host, err := validateHost(host)
	if err != nil {
		return nil, err
	}

	rp := &DummyProxy{
		host:    host,
		search:  search,
		replace: replace,
	}

	return rp, nil
}

type DummyProxy struct {
	host    string
	search  string
	replace string
}

func (rp *DummyProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(rp.host)
	if err != nil {
		log.Println("HTTP GET error:", err)
		return
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read body error:", err)
		return
	}

	b = bytes.Replace(b, []byte(rp.search), []byte(rp.replace), -1)

	w.Write(b)
}
