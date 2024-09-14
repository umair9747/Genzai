package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/rumble773/Genzai-UI/internal/models"
)

var (
	SaveOutput string
	NumWorkers int
	Timeout    int
	genzaiDB   models.GenzaiDB
	vendorDB   models.VendorDB
	vendorVulnsDB models.VendorVulnsDB
)

const banner = `
::::::::   :::::::::: ::::    ::: :::::::::     :::     ::::::::::: 
:+:    :+: :+:        :+:+:   :+:      :+:    :+: :+:       :+:     
+:+        +:+        :+:+:+  +:+     +:+    +:+   +:+      +:+     
:#:        +#++:++#   +#+ +:+ +#+    +#+    +#++:++#++:     +#+     
+#+   +#+# +#+        +#+  +#+#+#   +#+     +#+     +#+     +#+     
#+#    #+# #+#        #+#   #+#+#  #+#      #+#     #+#     #+#     
 ########  ########## ###    #### ######### ###     ### ########### 

	The IoT Security Toolkit by Umair Nehri (0x9747)
`

func PrintBanner() {
	fmt.Println(banner + "\n")
}

func LoadDB() models.GenzaiDB {
	fileContent, err := os.ReadFile("data/signatures.json")
	if err != nil {
		log.Println("error reading file: ", err)
		os.Exit(0)
	}

	var data map[string]models.GenzaiDB
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
	return genzaiDB
}

func LoadVendorDB() models.VendorDB {
	fileContent, err := ioutil.ReadFile("data/vendor-logins.json")
	if err != nil {
		log.Println("error reading file: ", err)
		os.Exit(0)
	}
	err = json.Unmarshal(fileContent, &vendorDB)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		os.Exit(0)
	}
	return vendorDB
}

func LoadVendorVulnsDB() models.VendorVulnsDB {
	fileContent, err := ioutil.ReadFile("data/vendor-vulns.json")
	if err != nil {
		log.Println("error reading file: ", err)
		os.Exit(0)
	}
	err = json.Unmarshal(fileContent, &vendorVulnsDB)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		os.Exit(0)
	}
	return vendorVulnsDB
}

func LogData(data string, filename string) {
	err := ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		log.Println("Error while logging the output", err.Error())
	}
	fmt.Printf("\nLogged the output in %s!", filename)
}