package proxy

import (
	"bufio"
	"bytes"
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanRunes)
	var buf bytes.Buffer
	for scanner.Scan() {
		_, err = buf.WriteString(scanner.Text())
		if err != nil {
			log.Println("Write to buffer error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	b := bytes.Replace(buf.Bytes(), []byte(rp.search), []byte(rp.replace), -1)

	_, err = w.Write(b)
	if err != nil {
		log.Println("Write to response error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
