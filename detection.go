package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func targetDetection(target string) string {
	resp, err := makeGetRequest(target)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return ""
	}
	defer resp.Body.Close()

	// PERFORM FINGERPRINTING

	// ITERATE OVER ALL DB ENTRIES
	for product, entry := range genzaiDB {

		// MATCH AGAINST HEADERS
		if entry.Matchers.Headers != nil {
			for headerKey, headerValue := range entry.Matchers.Headers {
				for key, values := range resp.Header {
					for _, value := range values {
						if headerKey == key && strings.Contains(value, headerValue.(string)) {
							return product
						}
					}
				}
			}
		}

		respBody, _ := ioutil.ReadAll(resp.Body)

		// MATCH AGAINST STRINGS
		for _, matchEntry := range entry.Matchers.Strings {
			if strings.Contains(string(respBody), matchEntry) {
				return product
			}
		}

		// MATCH AGAINST NON-200 CODES

		if entry.Matchers.ResponseCode != 200 {
			if entry.Matchers.ResponseCode == resp.StatusCode {
				return product
			}
		}
	}

	return ""
}

func makeGetRequest(url string) (*http.Response, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow redirects
			return nil
		},
		Timeout: 30 * time.Second, // Set timeout to 30 seconds
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
