package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	
	"github.com/rumble773/Genzai-UI/internal/detection"
	"github.com/rumble773/Genzai-UI/internal/models"
	"github.com/rumble773/Genzai-UI/internal/scan"
	"github.com/rumble773/Genzai-UI/internal/utils"
)

var (
	numWorkers int
	timeout    time.Duration
	verbose    bool
)

func main() {
	flag.StringVar(&utils.SaveOutput, "save", "", "Save the output in a file. [Default filename is output.json]")
	flag.IntVar(&numWorkers, "workers", 10, "Number of concurrent workers")
	flag.DurationVar(&timeout, "timeout", 30*time.Second, "Timeout for each scan")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/scan", handleScan).Methods("POST")
	r.HandleFunc("/health", handleHealth).Methods("GET")

	utils.PrintBanner()
	log.Println("Starting Genzai API server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Genzai API is healthy")
}

func handleScan(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var scanReq models.ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&scanReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	expandedTargets, err := expandTargets(scanReq.Targets)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error expanding targets: %v", err), http.StatusBadRequest)
		return
	}

	if len(expandedTargets) == 0 {
		http.Error(w, "No valid targets provided", http.StatusBadRequest)
		return
	}

	logVerbose("Expanded targets: %d", len(expandedTargets))

	log.Println("Genzai is starting...")
	log.Println("Loading Genzai Signatures DB...")
	genzaiDB := utils.LoadDB()
	log.Println("Loading Vendor Passwords DB...")
	vendorDB := utils.LoadVendorDB()
	log.Println("Loading Vendor Vulnerabilities DB...")
	vendorVulnsDB := utils.LoadVendorVulnsDB()

	results, errors := scanTargets(expandedTargets, genzaiDB, vendorDB, vendorVulnsDB)

	scanResponse := models.ScanResponse{
		Results:      results,
		Targets:      scanReq.Targets,
		TotalScanned: len(expandedTargets),
		TimeElapsed:  time.Since(start).String(),
		Errors:       errors,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scanResponse)
}

func expandTargets(targets []string) ([]string, error) {
	var expandedTargets []string
	for _, target := range targets {
		if ip := net.ParseIP(target); ip != nil {
			expandedTargets = append(expandedTargets, target)
		} else if _, ipNet, err := net.ParseCIDR(target); err == nil {
			for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
				expandedTargets = append(expandedTargets, ip.String())
			}
		} else if host, _, err := net.SplitHostPort(target); err == nil {
			expandedTargets = append(expandedTargets, host)
		} else {
			expandedTargets = append(expandedTargets, target)
		}
	}
	return expandedTargets, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func scanTargets(targets []string, genzaiDB models.GenzaiDB, vendorDB models.VendorDB, vendorVulnsDB models.VendorVulnsDB) ([]models.GenzaiResult, []string) {
	results := make([]models.GenzaiResult, 0, len(targets))
	errors := make([]string, 0)
	resultChan := make(chan models.GenzaiResult, len(targets))
	errorChan := make(chan string, len(targets))
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, numWorkers)

	for _, target := range targets {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(target string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			result, err := scanTarget(target, genzaiDB, vendorDB, vendorVulnsDB)
			if err != nil {
				errorChan <- fmt.Sprintf("Error scanning %s: %v", target, err)
			} else {
				resultChan <- result
			}
		}(target)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	for result := range resultChan {
		results = append(results, result)
	}

	for err := range errorChan {
		errors = append(errors, err)
	}

	return results, errors
}

func scanTarget(target string, genzaiDB models.GenzaiDB, vendorDB models.VendorDB, vendorVulnsDB models.VendorVulnsDB) (models.GenzaiResult, error) {
	logVerbose("Starting the scan for original target: %s", target)
	
	formattedTarget, err := ensureURLFormat(target)
	if err != nil {
		return models.GenzaiResult{}, fmt.Errorf("invalid target URL: %v", err)
	}
	
	logVerbose("Scanning formatted target: %s", formattedTarget)
	
	product, category, tag := detection.TargetDetection(formattedTarget, genzaiDB)
	if product == "" {
		logVerbose("No product identified for target: %s", formattedTarget)
		return models.GenzaiResult{
			Target: formattedTarget,
			Issues: []models.Issue{{
				IssueTitle: "No product identified",
				URL:        formattedTarget,
				AdditionalContext: "The target could not be identified as any known IoT product.",
			}},
		}, nil
	}

	result := models.GenzaiResult{
		Target:        formattedTarget,
		IoTidentified: product,
		Category:      category,
	}

	logVerbose("IoT Dashboard Discovered: %s", product)
	logVerbose("Trying for default vendor-specific [%s] passwords...", product)
	passIssue := scan.VendorpassScan(formattedTarget, product, tag, vendorDB)
	if passIssue.IssueTitle != "" {
		result.Issues = append(result.Issues, passIssue)
	}

	logVerbose("Scanning for any known vulnerabilities from the DB related to %s", product)
	vulnIssues := scan.VendorvulnScan(formattedTarget, product, tag, vendorVulnsDB)
	result.Issues = append(result.Issues, vulnIssues...)

	return result, nil
}

func ensureURLFormat(target string) (string, error) {
	target = strings.TrimSpace(target)

	// Check if it's an IP address
	if ip := net.ParseIP(target); ip != nil {
		return fmt.Sprintf("http://%s", target), nil
	}

	// If it doesn't have a scheme, add http://
	if !strings.Contains(target, "://") {
		target = "http://" + target
	}

	// Parse and validate the URL
	u, err := url.Parse(target)
	if err != nil {
		return "", err
	}

	// Ensure the scheme is either http or https
	if u.Scheme != "http" && u.Scheme != "https" {
		u.Scheme = "http"
	}

	return u.String(), nil
}

func logVerbose(format string, v ...interface{}) {
	if verbose {
		log.Printf(format, v...)
	}
}