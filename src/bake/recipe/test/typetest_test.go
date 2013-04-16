// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package test

import (
	"bytes"
	"strio"
	"testing"
)

func newLineReader(s string) strio.LineReader {
	buf := bytes.NewBufferString(s)
	return strio.NewLineReader(buf)
}

func TestValidDirective(t *testing.T) {
	// Arrange
	in := newLineReader("descr\n cmd")

	// Act
	tests, err := readTypeTests(in)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tests) != 1 {
		t.Fatalf("expected 1 parsed tests, got %d", len(tests))
	}

	if tests[0].descr() != "descr" {
		t.Errorf("expected 'descr', got '%s'", tests[0].descr())
	}

	if len(tests[0].actions()) != 1 {
		t.Errorf("expected 1 test action, got %d",
			len(tests[0].actions()))
	}
}

func TestEmptyDescr(t *testing.T) {
	// Arrange
	in := newLineReader("\n cmd")

	// Act
	_, err := readTypeTests(in)

	// Assert
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestInvalidDirective(t *testing.T) {
	// Arrange
	in := newLineReader("descr\n§cmd")

	// Act
	_, err := readTypeTests(in)

	// Assert
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestEOFAfterDescr(t *testing.T) {
	// Arrange
	in := newLineReader("descr")

	// Act
	_, err := readTypeTests(in)

	// Assert
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestEOFAfterDescrNewline(t *testing.T) {
	// Arrange
	in := newLineReader("descr\n")

	// Act
	_, err := readTypeTests(in)

	// Assert
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestExtraNewline(t *testing.T) {
	// Arrange
	in := newLineReader("descr\n cmd\n")

	// Act
	_, err := readTypeTests(in)

	// Assert
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestMultipleCommands(t *testing.T) {
	// Arrange
	in := newLineReader("descr\n cmd\n cmd\n cmd")

	// Act
	tests, err := readTypeTests(in)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tests) != 1 {
		t.Fatalf("expected 1 parsed tests, got %d", len(tests))
	}

	if tests[0].descr() != "descr" {
		t.Errorf("expected 'descr', got '%s'", tests[0].descr())
	}

	if len(tests[0].actions()) != 3 {
		t.Errorf("expected 3 test actions, got %d",
			len(tests[0].actions()))
	}
}

func TestMultipleTests(t *testing.T) {
	// Arrange
	in := newLineReader("descr1\n cmd\n\ndescr2\n cmd\n cmd")

	// Act
	tests, err := readTypeTests(in)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tests) != 2 {
		t.Fatalf("expected 2 parsed tests, got %d", len(tests))
	}

	if tests[0].descr() != "descr1" {
		t.Errorf("expected 'descr1', got '%s'", tests[0].descr())
	}

	if len(tests[0].actions()) != 1 {
		t.Errorf("expected 1 test actions, got %d",
			len(tests[0].actions()))
	}

	if tests[1].descr() != "descr2" {
		t.Errorf("expected 'descr2', got '%s'", tests[1].descr())
	}

	if len(tests[1].actions()) != 2 {
		t.Errorf("expected 2 test actions, got %d",
			len(tests[1].actions()))
	}
}

func TestEmptyTest(t *testing.T) {
	// Arrange
	in := newLineReader("descr1\n\ndescr2\n cmd\n cmd")

	// Act
	_, err := readTypeTests(in)

	// Assert
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}
