// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package readers

import (
	"fmt"
	"io"
)

type countingReader struct {
	in io.Reader

	// the number of characters read on the current line
	charNum int

	lineNum int
}

func NewCountingReader(in io.Reader) *countingReader {
	return &countingReader{in, 0, 1}
}

func (r *countingReader) Read(p []byte) (n int, err error) {
	n, err = r.in.Read(p)
	fmt.Printf("%d", n)

	for i := range p {
		if p[i] == '\n' {
			r.lineNum++
			r.charNum = 0
		} else {
			r.charNum++
		}
	}

	return
}

func (r *countingReader) CharNum() int {
	return r.charNum
}

func (r *countingReader) LineNum() int {
	return r.lineNum
}
