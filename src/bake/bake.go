// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package main provides the entry point to the bake executable.

package main

import (
	"bake/env"
	"bake/proj"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return "[" + strings.Join(*s, ", ") + "]"
}

func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

var (
	types stringSlice

	verbose = flag.Bool("v", false, "Print extra progress information")

	helpArgs = map[*bool]func(){
		flag.Bool("L", false, "Print supported languages"): printLangs,
	}

	optionalArgs = map[string]*string{
		"Email": flag.String("e", "", "Email address of the owner"),
	}

	lang  = flag.String("l", "", "Language of the project")
	owner = flag.String("o", "", "Owner of the project")
	name  = flag.String("n", "", "Name of the project")

	requiredArgs = map[string]*string{
		"l": lang,
		"o": owner,
		"n": name,
	}
)

func main() {
	flag.Var(&types, "t", "The project's types")

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

	vars := map[string]string{
		"ProjectName":      *name,
		"ProjectNameLower": strings.ToLower(*name),
		"Owner":            *owner,
	}

	for argName, argVal := range optionalArgs {
		if *argVal != "" {
			vars[argName] = *argVal
		}
	}

	p := proj.New(*lang, types, *verbose, vars)
	if err = p.GenTo(""); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "-%s is required\n", argName)
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

	langsSlice := sort.StringSlice(langs)
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
