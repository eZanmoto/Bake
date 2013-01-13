// Copyright {Year} {Owner}. All rights reserved.

// Package main provides the entry point to the {ProjectNameLower} executable.
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// vbose is 'true' when extra progress information output is encouraged. 
	vbose = flag.Bool("v", false, "Print extra progress information")

	// helpArgs contains options that output help pages if specified.
	helpArgs = map[*bool]func(){{
		flag.Bool("h", false, "Print usage information"): flag.Usage,
	}}

	// reqArgs contains options which require a value.
	reqArgs = map[string]*string{{
	}}
)

// main is the entry point to the {ProjectNameLower} executable.
func main() {{
	err := parseFlags()
	if err != nil {{
		fmt.Fprintf(os.Stderr, "%v\n", err)
		flag.Usage()
		os.Exit(1)
	}}

	fmt.Printf("{ProjectName} (C) {Year} {Owner}\n")
	if (*vbose) {{
		fmt.Printf("Run with -h to view usage information\n")
	}}
}}

// parseFlags parses the command-line arguments to the {ProjectNameLower} executable.
func parseFlags() {{
	flag.Parse()

	for argVal, printFunc := range helpArgs {{
		if *argVal {{
			printFunc()
			os.Exit(0)
		}}
	}}

	for argName, argVal := range reqArgs {{
		if *argVal == "" {{
			return fmt.Errorf("-%s is required", argName)
		}}
	}}
}}
