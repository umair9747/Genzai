package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func takeInput() {
	flag.Parse()
	args = flag.Args()

	if len(args) < 1 {
		fmt.Println("No arguments provied! [Exiting...]")
		os.Exit(0)
	} else {
		for _, arg := range args {
			if !strings.Contains(arg, ".") {

			}
			if strings.HasSuffix(arg, ".txt") {
				content, err := ioutil.ReadFile(arg)
				if err != nil {
					fmt.Println("Error reading file:", err)
					os.Exit(0)
				}
				targets = append(targets, strings.Split(string(content), "\n")...)
			}
			if !strings.HasPrefix(arg, "http") {
				targets = append(targets, "http://"+arg)
			} else {
				targets = append(targets, arg)
			}
		}
	}
}
