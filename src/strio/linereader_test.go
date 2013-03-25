// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package strio

import (
	"bytes"
	"io"
	"testing"
)

func newStringLineReader(str string) LineReader {
	in := bytes.NewBufferString(str)
	return NewLineReader(in)
}

func TestReadLine(t *testing.T) {
	// Arrange
	in := newStringLineReader("a\nb\nc")

	// Act
	str, err := in.ReadLine()

	// Assert
	if str != "a\n" {
		t.Errorf("Expected 'a\\n', got '%s'\n", str)
	}

	if err != nil {
		t.Errorf("Unexpected error '%v'\n", err)
	}
}

func TestReadEOFWithNewline(t *testing.T) {
	// Arrange
	in := newStringLineReader("a\nb\n")

	// Act
	in.ReadLine()
	b, err := in.ReadLine()
	blank, eof := in.ReadLine()

	// Assert
	if b != "b\n" {
		t.Errorf("Expected 'b\\n', got '%s'\n", b)
	}

	if err != nil {
		t.Errorf("Unexpected error '%v'\n", err)
	}

	if len(blank) != 0 {
		t.Errorf("Expected empty string, got '%s'\n", blank)
	}

	if eof != io.EOF {
		t.Errorf("Expected io.EOF, got '%v'\n", eof)
	}
}

func TestReadEOFWithoutNewline(t *testing.T) {
	// Arrange
	in := newStringLineReader("a\nb")

	// Act
	in.ReadLine()
	b, eof := in.ReadLine()

	// Assert
	if b != "b" {
		t.Errorf("Expected 'b', got '%s'\n", b)
	}

	if eof != io.EOF {
		t.Errorf("Expected io.EOF, got '%v'\n", eof)
	}
}

func TestChompLine(t *testing.T) {
	// Arrange
	in := newStringLineReader("a\nb\nc")

	// Act
	str, err := in.ChompLine()

	// Assert
	if str != "a" {
		t.Errorf("Expected 'a', got '%s'\n", str)
	}

	if err != nil {
		t.Errorf("Unexpected error '%v'\n", err)
	}
}

func TestChompEOFWithNewline(t *testing.T) {
	// Arrange
	in := newStringLineReader("a\nb\n")

	// Act
	in.ChompLine()
	b, err := in.ChompLine()
	blank, eof := in.ChompLine()

	// Assert
	if b != "b" {
		t.Errorf("Expected 'b', got '%s'\n", b)
	}

	if err != nil {
		t.Errorf("Unexpected error '%v'\n", err)
	}

	if len(blank) != 0 {
		t.Errorf("Expected empty string, got '%s'\n", blank)
	}

	if eof != io.EOF {
		t.Errorf("Expected io.EOF, got '%v'\n", eof)
	}
}

func TestChompEOFWithoutNewline(t *testing.T) {
	// Arrange
	in := newStringLineReader("a\nb")

	// Act
	in.ChompLine()
	b, eof := in.ChompLine()

	// Assert
	if b != "b" {
		t.Errorf("Expected 'b', got '%s'\n", b)
	}

	if eof != io.EOF {
		t.Errorf("Expected io.EOF, got '%v'\n", eof)
	}
}
