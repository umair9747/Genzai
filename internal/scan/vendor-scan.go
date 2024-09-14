package scan

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

    "github.com/rumble773/Genzai-UI/internal/models"
    "github.com/rumble773/Genzai-UI/internal/utils"
)

func vendorpassScan(target string, product string, tag string) Issue {
	var vendorpassIssue Issue
	if !strings.HasSuffix(target, "/") {
		target += "/"
	}
	for _, entry := range vendorDB.Entries {
		if entry.Tag == tag { // LOOK FOR THE EXACT VENDOR PASS ENTRY WE WANNA TRY
			for _, payloadPath := range entry.Payload.Paths {
				resp, err := makeHTTPRequest(target+payloadPath, entry.Payload.Headers, entry.Payload.Body, entry.Payload.Method)
				if err != nil { // IF THERE WAS AN ERROR MAKING THE REQ
					log.Println(err)
					return vendorpassIssue
				} else { // IF THERE WERE NO ERRORS

					//DO THE MATCHING OVER HERE
					// FIRST DO NON-200 STATUS CODE MATCHES HERE
					if entry.Matchers.ResponseCode != 200 {
						if resp.StatusCode == entry.Matchers.ResponseCode {
							log.Println(target, "[", product, "]", "is vulnerable with default password - ", entry.Issue)
							vendorpassIssue.IssueTitle = entry.Issue
							vendorpassIssue.URL = target + payloadPath
							vendorpassIssue.AdditionalContext = "The resulting non-200 status code matched with the one in DB."
							return vendorpassIssue
						}
					}

					// SECONDLY CHECK FOR THE RESPONSE PATH

					if entry.Matchers.Responsepath != "" {
						if strings.Contains(resp.Request.URL.Path, entry.Matchers.Responsepath) {
							log.Println(target, "[", product, "]", "is vulnerable with default password - ", entry.Issue)
							vendorpassIssue.IssueTitle = entry.Issue
							vendorpassIssue.URL = target + payloadPath
							vendorpassIssue.AdditionalContext = "The resulting URL path matched with the one in DB."
							return vendorpassIssue
						}
					}

					// THIRDLY CHECK OVER THE RESPONSE HEADERS
					if entry.Matchers.Headers != nil {
						for headerKey, headerValue := range entry.Matchers.Headers {
							for key, values := range resp.Header {
								for _, value := range values {
									if strings.EqualFold(strings.ToLower(headerKey), strings.ToLower(key)) && strings.Contains(strings.ToLower(value), strings.ToLower(headerValue)) {
										log.Println(target, "[", product, "]", "is vulnerable with default password - ", entry.Issue)
										vendorpassIssue.IssueTitle = entry.Issue
										vendorpassIssue.URL = target + payloadPath
										vendorpassIssue.AdditionalContext = "The resulting headers matched with those in the DB."
										return vendorpassIssue
									}
								}
							}
						}
					}

					respBody, _ := ioutil.ReadAll(resp.Body)
					// NEXT CHECK FOR STRINGS WITHIN RESPONSE BODY
					if entry.Matchers.Strings != nil {
						for _, matchString := range entry.Matchers.Strings {
							matchRe := regexp.MustCompile(strings.ToLower(matchString))
							if matchRe.MatchString(strings.ToLower(string(respBody))) {
								log.Println(target, "[", product, "]", "is vulnerable with default password - ", entry.Issue)
								vendorpassIssue.IssueTitle = entry.Issue
								vendorpassIssue.URL = target + payloadPath
								vendorpassIssue.AdditionalContext = "The resulting body had matching strings from the DB."
								return vendorpassIssue
							}
						}
					}
				}
				/*
					fmt.Printf("Vendor: %s\n", vendor)
					fmt.Println("Payload:")
					fmt.Printf("  Method: %s\n", entry.Payload.Method)
					fmt.Println("  Headers:")
					for key, value := range entry.Payload.Headers {
						fmt.Printf("    %s: %s\n", key, value)
					}

					fmt.Println("Matchers:")
					fmt.Printf("  Response Code: %d\n", entry.Matchers.ResponseCode)
					fmt.Println("  Strings:", entry.Matchers.Strings)
					fmt.Println("  Headers:")
					for key, value := range entry.Matchers.Headers {
						fmt.Printf("    %s: %s\n", key, value)
					}
				*/
			}
		}
	}
	return vendorpassIssue
}
