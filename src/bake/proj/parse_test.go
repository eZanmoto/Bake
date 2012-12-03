// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

import (
	"fmt"
	"testing"
	"tests/perm"
)

const (
	numTests = 1000
)

func TestParseInsert(t *testing.T) {
	letters := perm.NewStringPermuter("abc")
	numbers := perm.NewStringPermuter("012")

	for i := 0; i < numTests; i++ {
		before, after := letters.Permute(), numbers.Permute()
		p := &Project{"", false, map[string]string{before: after}}
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
		t.Fatalf("'%s': Expected '%s', got '%s'", before, after, result)
	}
}

func TestParseDoubleLBrace(t *testing.T) {
	p := &Project{"", false, map[string]string{}}
	letters := perm.NewStringPermuter("abc")

	for i := 0; i < numTests; i++ {
		padding := letters.Permute()
		expectParse(t, p, "{{"+padding, "{"+padding)
		expectParse(t, p, padding+"{{"+padding, padding+"{"+padding)
		expectParse(t, p, padding+"{{", padding+"{")
	}
}

func TestParseMissingClose(t *testing.T) {
	p := &Project{"", false, map[string]string{}}
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
	p := &Project{"", false, map[string]string{}}
	letters := perm.NewStringPermuter("abc")

	for i := 0; i < numTests; i++ {
		padding := letters.Permute()
		expectParse(t, p, "}}"+padding, "}"+padding)
		expectParse(t, p, padding+"}}"+padding, padding+"}"+padding)
		expectParse(t, p, padding+"}}", padding+"}")
	}
}

func TestParseSingleRBrace(t *testing.T) {
	p := &Project{"", false, map[string]string{}}
	letters := perm.NewStringPermuter("abc")

	for i := 0; i < numTests; i++ {
		parseFail(t, p, "}"+letters.Permute())
		parseFail(t, p, letters.Permute()+"}"+letters.Permute())
		parseFail(t, p, letters.Permute()+"}")
	}
}

func TestParseVarDepsInc(t *testing.T) {
	letters := perm.NewStringPermuter("abc")
	numbers := perm.NewStringPermuter("012")

	letters.Permute() // Skip empty string
	numbers.Permute() // Skip empty string

	for i := 0; i < numTests; i++ {
		z := letters.Permute()
		a, b := letters.Permute(), numbers.Permute()
		c, d := letters.Permute(), numbers.Permute()
		p := &Project{"", false, map[string]string{a: b, c: d}}

		expectParse(t, p,
			fmt.Sprintf("%s{?%s:{%s}}%s", z, a, a, z),
			z+b+z)

		expectParse(t, p,
			fmt.Sprintf("{?%s:%s{%s}%s}", z, a, a, z),
			"")

		expectParse(t, p,
			fmt.Sprintf("%s{?x:{%s}}%s", z, a, z),
			z+z)

		expectParse(t, p,
			fmt.Sprintf("%s{?x:{%s}%s}", z, a, z),
			z)

		expectParse(t, p,
			fmt.Sprintf("{?x:%s{%s}}%s", z, a, z),
			z)

		expectParse(t, p,
			fmt.Sprintf("{?x:%s{%s}%s}", z, a, z),
			"")

		expectParse(t, p,
			fmt.Sprintf("%s{?%s&%s:{%s}{%s}}%s", z, a, c, a, c, z),
			z+b+d+z)

		expectParse(t, p,
			fmt.Sprintf("%s{?x&%s:{x}{%s}}%s", z, a, a, z),
			z+z)

		expectParse(t, p,
			fmt.Sprintf("%s{?%s&x:{%s}{x}}%s", z, a, a, z),
			z+z)
	}
}
