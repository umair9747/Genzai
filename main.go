package main

import (
	"fmt"
	"log"
)

func main() {
	printBanner()
	takeInput()
	if len(targets) > 0 {
		log.Println("Genzai is starting...")
		log.Println("Loading Genzai Signatures DB...")
		loadDB()
		log.Println("Loading Vendor Passwords DB...")
		loadVendorDB()
		log.Println("Loading Vendor Vulnerabilities DB...")
		loadVendorVulnsDB()
		fmt.Println("\n ")
	}

	genzaiOutput.Targets = targets

	for _, target := range targets {
		fmt.Println()
		log.Println("Starting the scan for", target)
		product, category, tag := targetDetection(target)
		if product != "" {
			var targetResult genzaiResult
			targetResult.Target = target
			targetResult.IoTidentified = product
			targetResult.Category = category
			log.Println("IoT Dashboard Discovered:", product)
			log.Println("Trying for default vendor-specific [", product, "] passwords...")
			passIssue := vendorpassScan(target, product, tag)
			if passIssue.URL != "" {
				targetResult.Issues = append(targetResult.Issues, passIssue)
			}
			log.Println("Scanning for any known vulnerabilities from the DB related to", product)
			vulnIssues := vendorvulnScan(target, product, tag)
			for _, vulnIssue := range vulnIssues {
				if vulnIssue.URL != "" {
					targetResult.Issues = append(targetResult.Issues, vulnIssue)
				}
			}
			genzaiOutput.Results = append(genzaiOutput.Results, targetResult)
		}
	}

	/*
		for key, entry := range genzaiDB {
			fmt.Printf("Entry %s:\n", key)
			fmt.Println("Matches:", entry.Matchers.Strings)
			fmt.Println("Response Code:", entry.Matchers.ResponseCode)

			fmt.Println("Headers:")
			for headerKey, headerValue := range entry.Matchers.Headers {
				fmt.Printf("  %s: %v\n", headerKey, headerValue)
			}
			fmt.Println()
		}
	*/
	generateOutput()
}
