// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package template

import (
	"testing"
)

func TestExpandVar(t *testing.T) {
	// Arrange
	before, after := "a", "b"
	d := &Dict{before: after}

	// Act & Assert
	testExpand(t, d, before, before)
	testExpand(t, d, "{"+before+"}", after)
	testExpand(t, d, before+"{"+before+"}", before+after)
	testExpand(t, d, "{"+before+"}"+before, after+before)
	testExpand(t, d, before+"{"+before+"}"+before, before+after+before)
}

func testExpand(t *testing.T, d *Dict, before, after string) {
	result, err := d.ExpandStr(before)
	if err != nil {
		t.Fatalf("Unexpected error while expanding '%s': %v",
			before, err)
	} else if after != result {
		t.Fatalf("Expanding:\n%s\nExpected:\n%s\nGot:\n%s\n",
			before, after, result)
	}
}

func TestexpandEmptyVar(t *testing.T) {
	expandFail(t, &Dict{}, "{}")
}

func expandFail(t *testing.T, d *Dict, value string) {
	_, err := d.ExpandStr(value)
	if err == nil {
		t.Fatalf("Expected error while parsing '%s', got none", value)
	}
}

func TestExpandUnknownName(t *testing.T) {
	expandFail(t, &Dict{}, "{x}")
}

func TestExpandDoubleLBrace(t *testing.T) {
	padding := " "
	d := &Dict{}

	testExpand(t, d, "{{"+padding, "{"+padding)
	testExpand(t, d, padding+"{{"+padding, padding+"{"+padding)
	testExpand(t, d, padding+"{{", padding+"{")
}

func TestParseSingleRBrace(t *testing.T) {
	padding := " "
	d := &Dict{}

	expandFail(t, d, "}"+padding)
	expandFail(t, d, padding+"}"+padding)
	ExpandFail(t, d, padding+"}")
}
