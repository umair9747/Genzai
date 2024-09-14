package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func takeInput() {
	flag.StringVar(&saveOutput, "save", "", "Save the output in a file. [Default filename is output.json]")
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		fmt.Println("No arguments provied! [Exiting...]")
		os.Exit(0)
	} else {
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
