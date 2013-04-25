// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"bake/recipe"
	"fmt"
	"os"
	"time"
)

var (
	testLangs = []string{"go"}
)

func main() {
	t := time.Now()
	exitStatus := 0

	for _, lang := range testLangs {
		if passed, err := recipe.Test(lang); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			exitStatus = 1
		} else if !passed {
			exitStatus = 1
		}

		status := "ok"
		if exitStatus != 0 {
			status = "FAIL"
		}

		fmt.Printf("%s\t%s-recipe\t%.3fs\n",
			status, lang, time.Since(t).Seconds())
	}

	os.Exit(exitStatus)
}
