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

func main() {
	t := time.Now()

	exitStatus := 0
	if passed, err := recipe.Test("go"); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		exitStatus = 1
	} else if !passed {
		exitStatus = 1
	}

	status := "ok"
	if exitStatus != 0 {
		status = "FAIL"
	}
	fmt.Printf("%s\trcptest\t%.3fs\n", status, time.Since(t).Seconds())

	os.Exit(exitStatus)
}
