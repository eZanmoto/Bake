// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"bake/proj"
	"fmt"
	"os"
	"strings"
)

func usage() {
	fmt.Printf("%s incl-file\n", os.Args[0])
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(2)
	}

	fname := os.Args[1]
	node, err := proj.ParseInclsFile(fname)

	if err != nil {
		fmt.Printf("Error parsing '%s': %v", fname, err)
		os.Exit(2)
	}

	str := node.String()
	if len(str) > 2 {
		str := strings.Replace(str[3:], "\n\t", "\n", -1)
		fmt.Printf("%s\n", str)
	}
}
