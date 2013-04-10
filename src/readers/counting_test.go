// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package readers

import (
	"bytes"
	"testing"
)

// TODO pass to an io.Reader tester

func newStringCountingReader(str string) *countingReader {
	in := bytes.NewBufferString(str)
	return NewCountingReader(in)
}

func TestCountChars(t *testing.T) {
	// Arrange
	in := newStringCountingReader("abcde")
	p := make([]byte, 3)

	// Act
	n, err := in.Read(p)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}

	if str := string(p); str != "abc" {
		t.Errorf("expected to read 'abc', got '%s'", str)
	}

	if n != 3 {
		t.Errorf("expected to read 3 bytes, got %d", n)
	}

	if in.CharNum() != 3 {
		t.Errorf("expected to be at char 3, but at %d", in.CharNum())
	}

	if in.LineNum() != 1 {
		t.Errorf("expected to be at line 1, but at %d", in.LineNum())
	}
}

func TestCountLines(t *testing.T) {
	// Arrange
	in := newStringCountingReader("a\nb\nc\nd\ne")
	p := make([]byte, 4)

	// Act
	n, err := in.Read(p)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error '%v'", err)
	}

	if str := string(p); str != "a\nb\n" {
		t.Errorf("expected to read 'a\\nb\\n', got '%s'", str)
	}

	if n != 4 {
		t.Errorf("expected to read 4 bytes, got %d", n)
	}

	if in.CharNum() != 0 {
		t.Errorf("expected to be at char 0, but at %d", in.CharNum())
	}

	if in.LineNum() != 3 {
		t.Errorf("expected to be at line 3, but at %d", in.LineNum())
	}
}
