// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package template

import (
	"testing"
)

func TestExpandEmpty(t *testing.T) {
	testExpand(t, &Dict{}, "", "")
}

func TestExpandVar(t *testing.T) {
	before, after := "a", "b"
	d := &Dict{before: after}

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

	testExpand(t, d, "{{", "{")
	testExpand(t, d, "{{"+padding, "{"+padding)
	testExpand(t, d, padding+"{{", padding+"{")
	testExpand(t, d, padding+"{{"+padding, padding+"{"+padding)
}

func TestExpandSingleRBrace(t *testing.T) {
	padding := " "
	d := &Dict{}

	expandFail(t, d, "}")
	expandFail(t, d, "}"+padding)
	expandFail(t, d, padding+"}")
	expandFail(t, d, padding+"}"+padding)
}

func TestExpandDoubleRBrace(t *testing.T) {
	padding := " "
	d := &Dict{}

	testExpand(t, d, "}}", "}")
	testExpand(t, d, "}}"+padding, "}"+padding)
	testExpand(t, d, padding+"}}", padding+"}")
	testExpand(t, d, padding+"}}"+padding, padding+"}"+padding)
}

func TestExpandSimpleCond(t *testing.T) {
	d := &Dict{"x": "a", "y": ""}

	testExpand(t, d, "<{?x}{x}{?}>", "<a>")
	testExpand(t, d, "<{?x}c{x}{?}>", "<ca>")
	testExpand(t, d, "<{?x}{x}t{?}>", "<at>")
	testExpand(t, d, "<{?x}c{x}t{?}>", "<cat>")

	testExpand(t, d, "<{?z}{x}{?}>", "<>")
	testExpand(t, d, "<{?z}c{x}{?}>", "<>")
	testExpand(t, d, "<{?z}{x}t{?}>", "<>")
	testExpand(t, d, "<{?z}c{x}t{?}>", "<>")

	testExpand(t, d, "<{?x}x{?}>", "<x>")
	testExpand(t, d, "<{?y}y{?}>", "<y>")
	testExpand(t, d, "<{?z}z{?}>", "<>")

	testExpand(t, d, "<{?x}x{?y}y{?}{?}>", "<xy>")
	testExpand(t, d, "<{?x}x{?z}z{?}{?}>", "<x>")
	testExpand(t, d, "<{?z}z{?y}y{?}{?}>", "<>")
}

func TestExpandCondElse(t *testing.T) {
	d := &Dict{"x": "a", "y": ""}

	testExpand(t, d, "<{?x}{x}{:}b{?}>", "<a>")
	testExpand(t, d, "<{?z}{z}{:}b{?}>", "<b>")

	testExpand(t, d, "<{?x}1{:}2{?}>", "<1>")
	testExpand(t, d, "<{?z}1{:}2{?}>", "<2>")

	testExpand(t, d, "<{?x}{?y}1{:}2{?}{:}{?y}3{:}4{?}{?}>", "<1>")
	testExpand(t, d, "<{?x}{?z}1{:}2{?}{:}{?y}3{:}4{?}{?}>", "<2>")
	testExpand(t, d, "<{?z}{?z}1{:}2{?}{:}{?y}3{:}4{?}{?}>", "<3>")
	testExpand(t, d, "<{?z}{?z}1{:}2{?}{:}{?z}3{:}4{?}{?}>", "<4>")
}

func TestExpandCondElsif(t *testing.T) {
	d := &Dict{"x": "a"}

	testExpand(t, d, "<{?z}1{:z}2{:z}3{?}>", "<>")

	testExpand(t, d, "<{?x}1{:z}2{:z}3{:}4{?}>", "<1>")
	testExpand(t, d, "<{?z}1{:x}2{:z}3{:}4{?}>", "<2>")
	testExpand(t, d, "<{?z}1{:z}2{:x}3{:}4{?}>", "<3>")
	testExpand(t, d, "<{?z}1{:z}2{:z}3{:}4{?}>", "<4>")

	testExpand(t, d, "<{?x}{x}{:z}2{:z}3{:}4{?}>", "<a>")
	testExpand(t, d, "<{?z}1{:x}{x}{:z}3{:}4{?}>", "<a>")
	testExpand(t, d, "<{?z}1{:z}2{:x}{x}{:}4{?}>", "<a>")
	testExpand(t, d, "<{?z}1{:z}2{:z}3{:}{x}{?}>", "<a>")
}
