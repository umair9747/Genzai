package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func startAPIServer() {
	http.HandleFunc("/scan", apiHandler) // Handle POST requests at `/scan` endpoint
	log.Println("Starting API server on :8585...")
	log.Fatal(http.ListenAndServe(":8585", nil)) // Start server on port 8585
}

type ScanRequest struct {
	Targets []string `json:"targets"`
}

type ScanResponse struct {
	Output string `json:"output"`
}

// API handler that processes requests
func apiHandler(w http.ResponseWriter, r *http.Request) {
	// Set the custom response header
	w.Header().Set("genzai-api", "1.0")

	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Set the targets from the request body
	targets = req.Targets

	// Reset the output structure before starting a new scan
	genzaiOutput = Response{Targets: targets, Results: []genzaiResult{}}

	// Run the existing functionality
	log.Println("Genzai is starting...")
	log.Println("Loading Genzai Signatures DB...")
	loadDB()
	log.Println("Loading Vendor Passwords DB...")
	loadVendorDB()
	log.Println("Loading Vendor Vulnerabilities DB...")
	loadVendorVulnsDB()

	// Existing scanning logic
	for _, target := range targets {
		log.Println("Starting the scan for", target)
		product, category, tag := targetDetection(target)
		if product != "" {
			var targetResult genzaiResult
			targetResult.Target = target
			targetResult.IoTidentified = product
			targetResult.Category = category

			passIssue := vendorpassScan(target, product, tag)
			if passIssue.URL != "" {
				targetResult.Issues = append(targetResult.Issues, passIssue)
			}

			vulnIssues := vendorvulnScan(target, product, tag)
			for _, vulnIssue := range vulnIssues {
				if vulnIssue.URL != "" {
					targetResult.Issues = append(targetResult.Issues, vulnIssue)
				}
			}
			genzaiOutput.Results = append(genzaiOutput.Results, targetResult)
		}
	}

	// Generate JSON output
	genzaiJson, err := json.MarshalIndent(genzaiOutput, "", "    ")
	if err != nil {
		http.Error(w, "Failed to generate output", http.StatusInternalServerError)
		return
	}

	// Set content type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Set HTTP status code to 200 OK
	w.WriteHeader(http.StatusOK)
	// Write the JSON output
	w.Write(genzaiJson)
}
