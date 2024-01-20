package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func vendorpassScan(target string, product string) {
	if !strings.HasSuffix(target, "/") {
		target += "/"
	}
	for vendor, entry := range vendorDB.Entries {
		if vendor == product { // LOOK FOR THE EXACT VENDOR PASS ENTRY WE WANNA TRY
			for _, payloadPath := range entry.Payload.Paths {

				if entry.Payload.Method == "GET" { // CHECK THE PAYLOAD METHOD - IF ITS GET
					resp, err := makeHTTPRequestGET(target+payloadPath, entry.Payload.Headers)
					if err != nil { // IF THERE WAS AN ERROR MAKING THE REQ
						log.Println(err)
					} else { // IF THERE WERE NO ERRORS

						//DO THE MATCHING OVER HERE

						// FIRST DO NON-200 STATUS CODE MATCHES HERE
						if entry.Matchers.ResponseCode != 200 {
							if resp.StatusCode == entry.Matchers.ResponseCode {
								log.Fatalln("Discovered", entry.Issue)
							}
						}

						// SECONDLY CHECK OVER THE RESPONSE HEADERS
						if entry.Matchers.Headers != nil {
							for headerKey, headerValue := range entry.Matchers.Headers {
								for key, values := range resp.Header {
									for _, value := range values {
										if headerKey == key && strings.Contains(value, headerValue) {
											log.Fatalln("Discovered", entry.Issue)
										}
									}
								}
							}
						}

						respBody, _ := ioutil.ReadAll(resp.Body)
						// NEXT CHECK FOR STRINGS WITHIN RESPONSE BODY
						if entry.Matchers.Strings != nil {
							for _, matchString := range entry.Matchers.Strings {
								if strings.Contains(string(respBody), matchString) {
									log.Fatalln(target, "[", product, "]", "is vulnerable with default password - ", entry.Issue)
								}
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
				fmt.Println()
			}
		}
	}
}
