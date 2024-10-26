package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func targetDetection(target string) (string, string, string) {
	resp, err := makeGetRequest(target)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return "", "", ""
	}
	defer resp.Body.Close()

	// PERFORM FINGERPRINTING

	respBody, _ := io.ReadAll(resp.Body)
	// ITERATE OVER ALL DB ENTRIES
	for product, entry := range genzaiDB {
		if entry.Matchers.Condition == "OR" {
			// MATCH AGAINST HEADERS
			if entry.Matchers.Headers != nil {
				for headerKey, headerValue := range entry.Matchers.Headers {
					for key, values := range resp.Header {
						for _, value := range values {
							if strings.EqualFold(strings.ToLower(headerKey), strings.ToLower(key)) && strings.Contains(strings.ToLower(value), strings.ToLower(headerValue.(string))) {
								return product, entry.Category, entry.Tag
							}
						}
					}
				}
			}

			// MATCH AGAINST STRINGS
			for _, matchEntry := range entry.Matchers.Strings {
				if strings.Contains(string(respBody), matchEntry) {
					return product, entry.Category, entry.Tag
				}
			}

			// MATCH AGAINST NON-200 CODES

			if entry.Matchers.ResponseCode != 200 {
				if entry.Matchers.ResponseCode == resp.StatusCode {
					return product, entry.Category, entry.Tag
				}
			}
		} else if entry.Matchers.Condition == "AND" {
			if andConditionMatcher(entry, resp, respBody) {
				return product, entry.Category, entry.Tag
			}
		}
	}

	return "", "", ""
}

func makeGetRequest(url string) (*http.Response, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow redirects
			return nil
		},
		Timeout: 30 * time.Second, // Set timeout to 30 seconds
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disable SSL certificate verification
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func andConditionMatcher(entry Entry, resp *http.Response, respBody []byte) bool {
	var isMatched bool
	var headersMatched bool
	var stringsMatched bool
	var statusMatched bool

	if entry.Matchers.Headers != nil {
		var headerScore int
		for headerKey, headerValue := range entry.Matchers.Headers {
			for key, values := range resp.Header {
				for _, value := range values {
					if strings.EqualFold(strings.ToLower(headerKey), strings.ToLower(key)) && strings.Contains(strings.ToLower(value), strings.ToLower(headerValue.(string))) {
						headerScore++
					}
				}
			}
		}
		if headerScore == len(entry.Matchers.Headers) {
			headersMatched = true
		}
	} else {
		headersMatched = true
	}

	// MATCH AGAINST STRINGS

	if len(entry.Matchers.Strings) > 0 {
		var stringScore int
		for _, matchEntry := range entry.Matchers.Strings {
			if strings.Contains(string(respBody), matchEntry) {
				stringScore++
			}
		}

		if stringScore == len(entry.Matchers.Strings) {
			stringsMatched = true
		}
	} else {
		stringsMatched = true
	}

	// MATCH AGAINST NON-200 CODES

	if entry.Matchers.ResponseCode != 200 {
		if entry.Matchers.ResponseCode == resp.StatusCode {
			statusMatched = true
		}
	} else {
		statusMatched = true
	}

	if headersMatched && stringsMatched && statusMatched {
		isMatched = true
	}

	return isMatched
}
