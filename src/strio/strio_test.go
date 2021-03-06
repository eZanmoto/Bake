// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package strio

import (
	"bytes"
	"testing"
)

func TestChompLinesWithNewLine(t *testing.T) {
	testChompLines(t, "a\nb\nc\n", "a", "b", "c", "")
}

func TestChompLinesWithoutNewLine(t *testing.T) {
	testChompLines(t, "a\nb\nc", "a", "b", "c")
}

func testChompLines(t *testing.T, src string, exp ...string) {
	// Arrange
	in := bytes.NewBufferString(src)

	// Act
	lines, err := ChompLines(in)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(lines) != len(exp) {
		t.Errorf("expected %d lines, got %d", len(exp), len(lines))
	}

	for i, line := range exp {
		if lines[i] != line {
			t.Errorf("line %d should be '%s', got '%s'",
				i, line, lines[i])
		}
	}
}

func TestReadAllWithNewLine(t *testing.T) {
	testReadAll(t, "a\nb\nc\n")
}

func TestReadAllWithoutNewLine(t *testing.T) {
	testReadAll(t, "a\nb\nc")
}

func testReadAll(t *testing.T, src string) {
	// Arrange
	in := bytes.NewBufferString(src)

	// Act
	output, err := ReadAll(in)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != src {
		t.Errorf("expected '%s', got '%s'", src, output)
	}
}
