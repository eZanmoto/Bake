// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package main provides the entry point to the bake executable.

package main

import (
	"bake/env"
	"flag"
	"fmt"
	"os"
	"sort"
)

var (
	langs = flag.Bool("L", false, "Print supported languages")

	helpArgs = map[*bool]func(){
		langs: printLangs,
	}
)

var (
	lang  = flag.String("l", "", "The language of the project")
	owner = flag.String("o", "", "The owner of the project")
	name  = flag.String("n", "", "The name of the project")

	requiredArgs = map[string]*string{
		"l": lang,
		"o": owner,
		"n": name,
	}
)

func main() {
	parseFlags()

	exists, err := langExists(*lang)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(2)
	} else if !exists {
		fmt.Fprintf(os.Stderr, "'"+*lang+"' is not a valid language\n")
		fmt.Fprintf(os.Stderr, "Use -languages to see valid options\n")
		os.Exit(2)
	}
}

func parseFlags() {
	flag.Parse()

	for argVal, printFunc := range helpArgs {
		if *argVal {
			printFunc()
			os.Exit(0)
		}
	}

	for argName, argVal := range requiredArgs {
		if *argVal == "" {
			fmt.Fprintf(os.Stderr, "-"+argName+" is required\n")
			flag.Usage()
			os.Exit(2)
		}
	}
}

func langExists(lang string) (bool, error) {
	langs, err := env.SupportedLangs()

	if err != nil {
		return false, err
	}

	langsSlice := sort.StringSlice(langs[0:])
	langsSlice.Sort()
	return langsSlice.Search(lang) != langsSlice.Len(), nil
}

func printLangs() {
	langs, err := env.SupportedLangs()

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		return
	}

	for _, lang := range langs {
		fmt.Printf(lang + "\n")
	}
}