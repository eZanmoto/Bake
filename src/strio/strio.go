// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package strio provides methods for manipulating I/O using strings.
package strio

import (
	"io"
)

// Reads all lines from 'r', with trailing newlines removed
func ChompLines(r io.Reader) (lines []string, err error) {
	in := NewLineReader(r)
	lines = make([]string, 0)

	for err == nil {
		var line string
		line, err = in.ChompLine()

		lines = append(lines, line)
	}

	if err == io.EOF {
		err = nil
	}

	return
}
