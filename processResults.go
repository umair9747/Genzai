package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Issue represents an individual issue within the Results structure
type Issue struct {
	IssueTitle        string `json:"IssueTitle"`
	URL               string `json:"URL"`
	AdditionalContext string `json:"AdditionalContext"`
}

// Result represents each result item in the Results array
type genzaiResult struct {
	Target        string  `json:"Target"`
	IoTidentified string  `json:"IoTidentified"`
	Category      string  `json:"category"`
	Issues        []Issue `json:"Issues"`
}

// Response represents the overall structure of the JSON
type Response struct {
	Results []genzaiResult `json:"Results"`
	Targets []string       `json:"Targets"`
}

var genzaiOutput Response

func generateOutput() {
	genzaiJson, err := json.MarshalIndent(genzaiOutput, "", "    ")
	if err != nil {
		log.Println(genzaiOutput)
		log.Println(err)
	} else {
		if saveOutput == "" {
			fmt.Println("")
			log.Println("No file name detected to log the output. Skipping to printing it!")
		} else {
			logData(string(genzaiJson), saveOutput)
		}
		fmt.Println("\n ")
		fmt.Println(string(genzaiJson))
	}
}
