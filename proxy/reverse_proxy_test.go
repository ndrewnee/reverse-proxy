package proxy

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	replaceString = "THIS_STRING_SHOULD_BE_IN_BODY"
)

func TestNewReverseProxy(t *testing.T) {
	tests := []struct {
		host   string
		search string
		err    error
	}{
		{
			"google.com",
			"Google",
			ErrHostNotHttpOrHttps,
		},
		{
			"http://mytube.uz",
			"MyTube",
			nil,
		},
		{
			"https://stackoverflow.com",
			"Stack",
			nil,
		},
		{
			"https://medium.com",
			"Medium",
			nil,
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://localhost:3000", nil)

		reverseProxy, err := NewReverseProxy(tt.host, tt.search, replaceString)
		if err != tt.err {
			t.Log("Host:", tt.host)
			t.Log("Search:", tt.search)
			t.Errorf("Excpected %s, got %s\n", tt.err, err)
			continue
		}

		if reverseProxy == nil {
			continue
		}

		reverseProxy.ServeHTTP(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Log("Host:", tt.host)
			t.Log("Search:", tt.search)
			t.Error(err)
			continue
		}

		if !strings.Contains(string(body), replaceString) {
			t.Errorf("Text '%s' not found in body\n", replaceString)
			t.Log("Host:", tt.host)
			t.Log("Search:", tt.search)
			t.Logf("Body: %s\n", body)
		}
	}
}
