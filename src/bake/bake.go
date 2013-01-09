// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package main provides the entry point to the bake executable.

package main

import (
	"bake/env"
	"bake/proj"
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return "[" + strings.Join(*s, ", ") + "]"
}

func (s *stringSlice) Set(vs string) error {
	for _, v := range strings.Split(vs, ",") {
		*s = append(*s, v)
	}
	return nil
}

var (
	types stringSlice

	verbose   = flag.Bool("v", false, "Print extra progress information")
	langTypes = flag.String("T", "", "Print project types for language")

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

	validateLang(*lang)

	vars := makeProjVars()
	for argName, argVal := range optionalArgs {
		if *argVal != "" {
			vars[argName] = *argVal
		}
	}

	p := proj.New(*lang, types, *verbose, vars)
	if err := p.GenTo(""); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(2)
	}
}

func makeProjVars() map[string]string {
	return map[string]string{
		"ProjectName":      *name,
		"ProjectNameLower": strings.ToLower(*name),
		"Owner":            *owner,
		"Year":             strconv.Itoa(time.Now().Year()),
	}
}

func parseFlags() {
	flag.Parse()

	if *langTypes != "" {
		validateLang(*langTypes)
		printTypesFor(*langTypes)
		os.Exit(0)
	}

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

func validateLang(lang string) {
	langs, err := env.SupportedLangs()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	langsSlice := sort.StringSlice(langs)
	langsSlice.Sort()

	if langsSlice.Search(lang) == langsSlice.Len() {
		fmt.Fprintf(os.Stderr, "'%s' is not a valid language\n", lang)
		fmt.Fprintf(os.Stderr, "Use -languages to see valid options\n")
		os.Exit(2)
	}
}

func printTypesFor(lang string) {
	tmplDir, err := env.TemplatesPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	langDir := path.Join(tmplDir, lang)

	fis, err := ioutil.ReadDir(langDir)
	if len(fis) <= 2 {
		fmt.Fprintf(os.Stderr, "'%s' is not fully supported\n", lang)
		os.Exit(2)
	}

	typeNames := make([]string, len(fis)-2)
	i := 0
	for _, fi := range fis {
		name := fi.Name()
		if name != "base" && name != "{ProjectName}" {
			typeNames[i] = name
			i++
		}
	}

	for _, name := range typeNames {
		file, err := os.Open(path.Join(langDir, name))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(2)
		}
		defer file.Close()

		descr, err := bufio.NewReader(file).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(2)
		}

		fmt.Printf("%s\t%s", name, descr)
	}
}

func printLangs() {
	langs, err := env.SupportedLangs()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	for _, lang := range langs {
		fmt.Printf("%s\n", lang)
	}
}
