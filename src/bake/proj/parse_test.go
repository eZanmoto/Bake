// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

import (
	"testing"
	"tests/perm"
)

const (
	numTests = 1000
)

func TestParseInsert(t *testing.T) {
	letters := perm.NewStringPermuter("abc")
	numbers := perm.NewStringPermuter("012")

	letters.Permute() // skip empty string
	numbers.Permute() // skip empty string

	for i := 0; i < numTests; i++ {
		before, after := letters.Permute(), numbers.Permute()
		p := &Project{vars: map[string]string{before: after}}
		expectParse(t, p, before, before)
		expectParse(t, p, "{"+before+"}", after)
		expectParse(t, p, before+"{"+before+"}", before+after)
		expectParse(t, p, "{"+before+"}"+before, after+before)
	}
}

func expectParse(t *testing.T, p *Project, before, after string) {
	result, err := p.parseStr(before)
	if err != nil {
		t.Fatalf("Unexpected error while parsing '%s': %v", before, err)
	} else if after != result {
		t.Fatalf("Parsing:\n%s\nExpected:\n%s\nGot:\n%s\n",
			before, after, result)
	}
}

func TestParseEmptyVar(t *testing.T) {
	parseFail(t, &Project{}, "{}")
}

func TestParseUnknownName(t *testing.T) {
	parseFail(t, &Project{}, "{X}")
}

func TestParseDoubleLBrace(t *testing.T) {
	p := &Project{}
	letters := perm.NewStringPermuter("abc")

	for i := 0; i < numTests; i++ {
		padding := letters.Permute()
		expectParse(t, p, "{{"+padding, "{"+padding)
		expectParse(t, p, padding+"{{"+padding, padding+"{"+padding)
		expectParse(t, p, padding+"{{", padding+"{")
	}
}

func TestParseMissingClose(t *testing.T) {
	p := &Project{}
	letters := perm.NewStringPermuter("abc")

	for i := 0; i < numTests; i++ {
		parseFail(t, p, "{"+letters.Permute())
		parseFail(t, p, letters.Permute()+"{"+letters.Permute())
		parseFail(t, p, letters.Permute()+"{")
	}
}

func parseFail(t *testing.T, p *Project, value string) {
	_, err := p.parseStr(value)
	if err == nil {
		t.Fatalf("Expected error while parsing '%s', got none", value)
	}
}

func TestParseDoubleRBrace(t *testing.T) {
	p := &Project{}
	letters := perm.NewStringPermuter("abc")

	for i := 0; i < numTests; i++ {
		padding := letters.Permute()
		expectParse(t, p, "}}"+padding, "}"+padding)
		expectParse(t, p, padding+"}}"+padding, padding+"}"+padding)
		expectParse(t, p, padding+"}}", padding+"}")
	}
}

func TestParseSingleRBrace(t *testing.T) {
	p := &Project{}
	letters := perm.NewStringPermuter("abc")

	for i := 0; i < numTests; i++ {
		parseFail(t, p, "}"+letters.Permute())
		parseFail(t, p, letters.Permute()+"}"+letters.Permute())
		parseFail(t, p, letters.Permute()+"}")
	}
}

func TestParseVarDepsInc(t *testing.T) {
	p := &Project{vars: map[string]string{"Y": "0", "Z": "2"}}

	expectParse(t, p, "a{?Y}{Y}{!}b", "a0b")
	expectParse(t, p, "a\n{?Y}{Y}{!}b", "a\n0b")
	expectParse(t, p, "a\n{?Y}\n{Y}{!}b", "a\n0b")
	expectParse(t, p, "a{?Y}\n{Y}{!}b", "a0b")
	expectParse(t, p, "a\n{?Y}\n{Y}\n{!}b", "a\n0\nb")
	expectParse(t, p, "a\n{?Y}\n{Y}\n{!}\nb", "a\n0\nb")
	expectParse(t, p, "a\n{?Y}\n{Y}\n{!}\n\nb", "a\n0\n\nb")

	expectParse(t, p, "a{?X}{X}{!}b", "ab")
	expectParse(t, p, "a\n{?X}{X}{!}b", "a\nb")
	expectParse(t, p, "a\n{?X}\n{X}{!}b", "a\nb")
	expectParse(t, p, "a{?X}\n{X}{!}b", "ab")
	expectParse(t, p, "a\n{?X}\n{X}\n{!}b", "a\nb")
	expectParse(t, p, "a\n{?X}\n{X}\n{!}\nb", "a\nb")
	expectParse(t, p, "a\n{?X}\n{X}\n{!}\n\nb", "a\n\nb")

	/*
		expectParse(t, p, "a{?Y}{Y}{?Z}{Z}{!}{!}b", "a02b")
		expectParse(t, p, "a{?X}0{:Y}1{:}2{!}b", "a1b")
		expectParse(t, p, "a{?X}0{:X}1{:}2{!}b", "a2b")
		expectParse(t, p, "a{?X}0{:X}1{:X}2{!}b", "ab")
	*/
}

func TestParseTypeDepsInc(t *testing.T) {
	p := &Project{types: []string{"y", "z"}}

	expectParse(t, p, "a{?y}y{!}b", "ayb")
	expectParse(t, p, "a\n{?y}y{!}b", "a\nyb")
	expectParse(t, p, "a\n{?y}\ny{!}b", "a\nyb")
	expectParse(t, p, "a{?y}\ny{!}b", "ayb")
	expectParse(t, p, "a\n{?y}\ny\n{!}b", "a\ny\nb")
	expectParse(t, p, "a\n{?y}\ny\n{!}\nb", "a\ny\nb")
	expectParse(t, p, "a\n{?y}\ny\n{!}\n\nb", "a\ny\n\nb")

	expectParse(t, p, "a{?x}x{!}b", "ab")
	expectParse(t, p, "a\n{?x}x{!}b", "a\nb")
	expectParse(t, p, "a\n{?x}\nx{!}b", "a\nb")
	expectParse(t, p, "a{?x}\nx{!}b", "ab")
	expectParse(t, p, "a\n{?x}\nx\n{!}b", "a\nb")
	expectParse(t, p, "a\n{?x}\nx\n{!}\nb", "a\nb")
	expectParse(t, p, "a\n{?x}\nx\n{!}\n\nb", "a\n\nb")

	/*
		expectParse(t, p, "a{?y}y{?z}z{?}{?}b", "ayzb")
		expectParse(t, p, "a{?x&x}{x}{x}{?}b", "ayzb")
		expectParse(t, p, "a{?y&x}{y}{x}{?}b", "ayzb")
		expectParse(t, p, "a{?x&z}{x}{z}{?}b", "ayzb")
		expectParse(t, p, "a{?y&z}{y}{z}{?}b", "ayzb")
		expectParse(t, p, "a{?x&x}{x}{x}{?}b", "ayzb")
		expectParse(t, p, "a{?y&x}{y}{x}{?}b", "ayzb")
		expectParse(t, p, "a{?x&z}{x}{z}{?}b", "ayzb")
		expectParse(t, p, "a{?y&z}{y}{z}{?}b", "ayzb")
		expectParse(t, p, "a{?x}0{:y}1{:}2{?}b", "a1b")
		expectParse(t, p, "a{?x}0{:x}1{:}2{?}b", "a2b")
		expectParse(t, p, "a{?x}0{:x}1{:x}2{?}b", "ab")
	*/
}
