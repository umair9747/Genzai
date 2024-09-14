package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/rumble773/Genzai-UI/internal/models"
)

var genzaiOutput models.Response

func GenerateOutput() {
	genzaiJson, err := json.MarshalIndent(genzaiOutput, "", "    ")
	if err != nil {
		log.Println(genzaiOutput)
		log.Println(err)
	} else {
		if SaveOutput == "" {
			fmt.Println("")
			log.Println("No file name detected to log the output. Skipping to printing it!")
		} else {
			LogData(string(genzaiJson), SaveOutput)
		}
		fmt.Println("\n ")
		fmt.Println(string(genzaiJson))
	}
}