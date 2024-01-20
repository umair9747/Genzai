package main

import (
	"net/http"
	"time"
)

func makeHTTPRequestGET(url string, headers map[string]string) (*http.Response, error) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	var req *http.Request
	var err error

	req, err = http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
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
