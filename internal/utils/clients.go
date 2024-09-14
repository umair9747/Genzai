package utils

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func MakeHTTPRequest(url string, headers map[string]string, body string, method string) (*http.Response, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var req *http.Request
	var err error

	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(http.MethodGet, url, nil)
	case http.MethodPost:
		if body != "" {
			req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(body))
		} else {
			req, err = http.NewRequest(http.MethodPost, url, nil)
		}
	default:
		return nil, fmt.Errorf("unsupported request method encountered")
	}

	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return client.Do(req)
}