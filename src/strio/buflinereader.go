// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package strio

import (
	"bufio"
	"io"
	"os"
)

type bufLineReader struct {
	in *bufio.Reader
}

func NewFileLineReader(path string) (LineReader, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return NewLineReader(file), nil
}

func NewLineReader(in io.Reader) LineReader {
	bufin := bufio.NewReader(in)
	return &bufLineReader{bufin}
}

func (r *bufLineReader) ReadLine() (string, error) {
	return r.in.ReadString('\n')
}

func (r *bufLineReader) ChompLine() (string, error) {
	s, e := r.ReadLine()

	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}

	return s, e
}
