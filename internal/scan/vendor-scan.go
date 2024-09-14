package scan

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/rumble773/Genzai-UI/internal/models"
	"github.com/rumble773/Genzai-UI/internal/utils"
)

func VendorpassScan(target string, product string, tag string, vendorDB models.VendorDB) models.Issue {
	var vendorpassIssue models.Issue
	if !strings.HasSuffix(target, "/") {
		target += "/"
	}
	for _, entry := range vendorDB {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		if entryMap["Tag"].(string) == tag {
			payload := entryMap["Payload"].(map[string]interface{})
			for _, pathInterface := range payload["Paths"].([]interface{}) {
				payloadPath := pathInterface.(string)
				headers := make(map[string]string)
				for k, v := range payload["Headers"].(map[string]interface{}) {
					headers[k] = v.(string)
				}
				resp, err := utils.MakeHTTPRequest(target+payloadPath, headers, payload["Body"].(string), payload["Method"].(string))
				if err != nil {
					log.Println(err)
					return vendorpassIssue
				}

				matchers := entryMap["Matchers"].(map[string]interface{})

				// Check response code
				if matchers["ResponseCode"].(float64) != 200 {
					if int(matchers["ResponseCode"].(float64)) == resp.StatusCode {
						log.Println(target, "[", product, "]", "is vulnerable with default password - ", entryMap["Issue"])
						vendorpassIssue.IssueTitle = entryMap["Issue"].(string)
						vendorpassIssue.URL = target + payloadPath
						vendorpassIssue.AdditionalContext = "The resulting non-200 status code matched with the one in DB."
						return vendorpassIssue
					}
				}

				// Check response path
				if responsePath, ok := matchers["Responsepath"].(string); ok && responsePath != "" {
					if strings.Contains(resp.Request.URL.Path, responsePath) {
						log.Println(target, "[", product, "]", "is vulnerable with default password - ", entryMap["Issue"])
						vendorpassIssue.IssueTitle = entryMap["Issue"].(string)
						vendorpassIssue.URL = target + payloadPath
						vendorpassIssue.AdditionalContext = "The resulting URL path matched with the one in DB."
						return vendorpassIssue
					}
				}

				// Check response headers
				if matcherHeaders, ok := matchers["Headers"].(map[string]interface{}); ok {
					for headerKey, headerValue := range matcherHeaders {
						for key, values := range resp.Header {
							for _, value := range values {
								if strings.EqualFold(strings.ToLower(headerKey), strings.ToLower(key)) && strings.Contains(strings.ToLower(value), strings.ToLower(headerValue.(string))) {
									log.Println(target, "[", product, "]", "is vulnerable with default password - ", entryMap["Issue"])
									vendorpassIssue.IssueTitle = entryMap["Issue"].(string)
									vendorpassIssue.URL = target + payloadPath
									vendorpassIssue.AdditionalContext = "The resulting headers matched with those in the DB."
									return vendorpassIssue
								}
							}
						}
					}
				}

				respBody, _ := ioutil.ReadAll(resp.Body)
				// Check for strings within response body
				if matcherStrings, ok := matchers["Strings"].([]interface{}); ok {
					for _, matchString := range matcherStrings {
						matchRe := regexp.MustCompile(strings.ToLower(matchString.(string)))
						if matchRe.MatchString(strings.ToLower(string(respBody))) {
							log.Println(target, "[", product, "]", "is vulnerable with default password - ", entryMap["Issue"])
							vendorpassIssue.IssueTitle = entryMap["Issue"].(string)
							vendorpassIssue.URL = target + payloadPath
							vendorpassIssue.AdditionalContext = "The resulting body had matching strings from the DB."
							return vendorpassIssue
						}
					}
				}
			}
		}
	}
	return vendorpassIssue
}