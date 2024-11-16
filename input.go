package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var apiMode bool       // Global flag for API mode
var reconMode bool     // Global flag for Recon mode
var reconPorts []int   // Array of ports to scan in Recon mode
var reconSubnet string // Variable for subnet in Recon mode

func takeInput() {
	// Define `--api`, `--recon`, `--reconports`, and `--subnet` flags
	flag.BoolVar(&apiMode, "api", false, "Run the tool in API mode")
	flag.BoolVar(&reconMode, "recon", false, "Run the tool in Recon mode against local network")
	flag.StringVar(&saveOutput, "save", "", "Save the output in a file. [Default filename is output.json]")
	flag.StringVar(&reconSubnet, "subnet", "192.168.1.", "Specify the subnet for recon mode (e.g., 192.168.1., 10.0.0.)")
	flag.Func("reconports", "List of ports to scan for each active host in recon mode", func(s string) error {
		for _, p := range strings.Split(s, ",") {
			port, err := strconv.Atoi(p)
			if err != nil {
				return err
			}
			reconPorts = append(reconPorts, port)
		}
		return nil
	})
	flag.Parse()
	args = flag.Args()

	if apiMode {
		startAPIServer() // Start the API server if `--api` is passed
		return
	}

	if reconMode {
		log.Println("Recon mode activated. Scanning the local network...")
		runReconMode() // Call recon mode function
		return
	}

	if len(args) < 1 {
		fmt.Println("No arguments provided! [Exiting...]")
		os.Exit(0)
	} else {
		// Existing input handling logic
		for i := 0; i < len(args); i++ {
			arg := args[i]
			if arg == "save" || arg == "-save" || arg == "--save" {
				if i+1 < len(args) {
					if strings.HasSuffix(args[i+1], ".txt") || strings.HasSuffix(args[i+1], ".json") {
						saveOutput = args[i+1]
					} else {
						saveOutput = "output.json"
					}
				} else {
					saveOutput = "output.json"
				}
				i++ // Skip the next argument since it has been processed
			} else if strings.Contains(arg, ".") {
				if strings.HasSuffix(arg, ".txt") {
					content, err := ioutil.ReadFile(arg)
					if err != nil {
						fmt.Println("Error reading file:", err)
						os.Exit(0)
					}
					targs := strings.Split(string(content), "\n")
					for _, targ := range targs {
						if targ != "" {
							targets = append(targets, targ)
						}
					}
				} else {
					if !strings.HasPrefix(arg, "http") {
						targets = append(targets, "http://"+arg)
					} else {
						targets = append(targets, arg)
					}
				}
			}
		}
	}
}

func runReconMode() {
	activeHosts := scanLocalNetwork(reconSubnet)
	for _, host := range activeHosts {
		if len(reconPorts) > 0 {
			// Append ports if reconPorts flag is set
			for _, port := range reconPorts {
				target := fmt.Sprintf("http://%s:%d", host, port)
				targets = append(targets, target)
				log.Printf("Added host %s with port %d to targets", host, port)
			}
		} else {
			// Just add the host if no ports specified
			targets = append(targets, "http://"+host)
			log.Printf("Added host %s to targets", host)
		}
	}
	log.Println("Recon mode scan complete. Found targets:", targets)
}

func scanLocalNetwork(subnet string) []string {
	var activeHosts []string
	var mu sync.Mutex // To synchronize access to activeHosts
	concurrency := 20 // Adjust this number based on system capabilities
	sem := make(chan struct{}, concurrency)
	results := make(chan string, 254)

	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s%d", subnet, i)
		sem <- struct{}{} // Acquire a slot in the semaphore
		go func(ip string) {
			defer func() { <-sem }() // Release the slot in the semaphore
			if pingHost(ip) {
				log.Printf("Active host found: %s\n", ip)
				results <- ip
			}
		}(ip)
	}

	// Close results channel when all goroutines are done
	go func() {
		for i := 0; i < cap(sem); i++ {
			sem <- struct{}{}
		}
		close(results)
	}()

	// Collect results from the channel
	for ip := range results {
		mu.Lock()
		activeHosts = append(activeHosts, ip)
		mu.Unlock()
	}

	return activeHosts
}

func pingHost(ip string) bool {
	timeout := 1 * time.Second
	cmd := exec.Command("ping", "-c", "1", "-W", fmt.Sprintf("%d", int(timeout.Seconds())), ip)

	if err := cmd.Run(); err != nil {
		return false // Host is not reachable
	}
	return true // Host is reachable
}
