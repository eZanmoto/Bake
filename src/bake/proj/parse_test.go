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

	for i := 0; i < numTests; i++ {
		before, after := letters.Permute(), numbers.Permute()
		p := &Project{"", false, map[string]string{before: after}}

		result, err := p.parseStr(before)
		if err != nil {
			t.Fatalf("Unexpected error while parsing '%s': %v",
				before, err)
		} else if result != before {
			t.Errorf("Expected '%s', got '%s'", before, result)
		}

		before = "{" + before + "}"
		result, err = p.parseStr(before)
		if err != nil {
			t.Fatalf("Unexpected error while parsing '%s': %v",
				before, err)
		} else if result != after {
			t.Errorf("Expected '%s', got '%s'", after, result)
		}
	}
}

func TestParseDoubleLBrace(t *testing.T) {
	p := &Project{"", false, map[string]string{}}
	letters := perm.NewStringPermuter("abc")

	for i := 0; i < numTests; i++ {
		padding := letters.Permute()
		s := "{{" + padding
		expected := "{" + padding
		result, err := p.parseStr(s)
		if err != nil {
			t.Fatalf("Unexpected error parsing '%s': %v", s, err)
		} else if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}

		s = padding + "{{" + padding
		expected = padding + "{" + padding
		result, err = p.parseStr(s)
		if err != nil {
			t.Fatalf("Unexpected error parsing '%s': %v", s, err)
		} else if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}

		s = padding + "{{"
		expected = padding + "{"
		result, err = p.parseStr(s)
		if err != nil {
			t.Fatalf("Unexpected error parsing '%s': %v", s, err)
		} else if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	}
}

func TestParseMissingClose(t *testing.T) {
	p := &Project{"", false, map[string]string{}}
	letters := perm.NewStringPermuter("abc")

	for i := 0; i < numTests; i++ {
		s := "{" + letters.Permute()
		_, err := p.parseStr(s)
		if err == nil {
			t.Fatalf("Expected error parsing '%s', got none", s)
		}

		s = letters.Permute() + "{" + letters.Permute()
		_, err = p.parseStr(s)
		if err == nil {
			t.Fatalf("Expected error parsing '%s', got none", s)
		}

		s = letters.Permute() + "{"
		_, err = p.parseStr(s)
		if err == nil {
			t.Fatalf("Expected error parsing '%s', got none", s)
		}
	}
}

func TestParseDoubleRBrace(t *testing.T) {
	p := &Project{"", false, map[string]string{}}
	letters := perm.NewStringPermuter("abc")

	for i := 0; i < numTests; i++ {
		padding := letters.Permute()
		s := "}}" + padding
		expected := "}" + padding
		result, err := p.parseStr(s)
		if err != nil {
			t.Fatalf("Unexpected error parsing '%s': %v", s, err)
		} else if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}

		s = padding + "}}" + padding
		expected = padding + "}" + padding
		result, err = p.parseStr(s)
		if err != nil {
			t.Fatalf("Unexpected error parsing '%s': %v", s, err)
		} else if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}

		s = padding + "}}"
		expected = padding + "}"
		result, err = p.parseStr(s)
		if err != nil {
			t.Fatalf("Unexpected error parsing '%s': %v", s, err)
		} else if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	}
}

func TestParseSingleRBrace(t *testing.T) {
	p := &Project{"", false, map[string]string{}}
	letters := perm.NewStringPermuter("abc")

	for i := 0; i < numTests; i++ {
		s := "}" + letters.Permute()
		_, err := p.parseStr(s)
		if err == nil {
			t.Errorf("Expected error parsing '%s', got none", s)
		}

		s = letters.Permute() + "}" + letters.Permute()
		_, err = p.parseStr(s)
		if err == nil {
			t.Errorf("Expected error parsing '%s', got none", s)
		}

		s = letters.Permute() + "}"
		_, err = p.parseStr(s)
		if err == nil {
			t.Errorf("Expected error parsing '%s', got none", s)
		}
	}
}
