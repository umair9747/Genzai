package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func printBanner() {
	fmt.Println(banner + "\n")
}

func loadDB() {
	fileContent, err := os.ReadFile("signatures.json")
	if err != nil {
		log.Println("error reading file: ", err)
		os.Exit(0)
	}

	var data map[string]DynamicEntries
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		os.Exit(0)
	}

	var ok bool
	genzaiDB, ok = data["entries"]
	if !ok {
		log.Println("Invalid JSON format: missing 'entries'")
		os.Exit(0)
	}
}

func loadVendorDB() {
	fileContent, err := ioutil.ReadFile("vendor-logins.json")
	if err != nil {
		log.Println("error reading file: ", err)
		os.Exit(0)
	}
	err = json.Unmarshal(fileContent, &vendorDB)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		os.Exit(0)
	}
}

func loadVendorVulnsDB() {
	fileContent, err := ioutil.ReadFile("vendor-vulns.json")
	if err != nil {
		log.Println("error reading file: ", err)
		os.Exit(0)
	}
	err = json.Unmarshal(fileContent, &vendorVulnsDB)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		os.Exit(0)
	}
}
func logData(data string, filename string) {
	// Write the string data to the file
	err := ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		log.Println("Error while logging the output", err.Error())
	}
	fmt.Printf("\nLogged the output in %s!", filename)
}
