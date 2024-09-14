package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func makeHTTPRequest(url string, headers map[string]string, body string, method string) (*http.Response, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disable SSL certificate verification
		},
	}

	var req *http.Request
	var err error

	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
	case http.MethodPost:
		if body != "" {
			req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(body))
			if err != nil {
				return nil, err
			}
		} else {
			req, err = http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("unsupported request method encountered")
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
