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
	fileContent, err := ioutil.ReadFile("signatures.json")
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
