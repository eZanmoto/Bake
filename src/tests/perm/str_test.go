// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package perm

import "testing"

func TestFirstStrPerms(t *testing.T) {
	p := NewStringPermuter("a0")

	expectNextStr(t, p, "")
	expectNextStr(t, p, "a")
	expectNextStr(t, p, "0")
	expectNextStr(t, p, "aa")
	expectNextStr(t, p, "0a")
	expectNextStr(t, p, "a0")
	expectNextStr(t, p, "00")
	expectNextStr(t, p, "aaa")
	expectNextStr(t, p, "0aa")
	expectNextStr(t, p, "a0a")
	expectNextStr(t, p, "00a")
	expectNextStr(t, p, "aa0")
	expectNextStr(t, p, "0a0")
	expectNextStr(t, p, "a00")
	expectNextStr(t, p, "000")
	expectNextStr(t, p, "aaaa")
}

func expectNextStr(t *testing.T, p *StringPermuter, s string) {
	val := p.Permute()

	if len(val) != len(s) {
		t.Fatalf("Expected length of %d, got %d", len(s), len(val))
	}

	for i := 0; i < len(val); i++ {
		if val[i] != s[i] {
			t.Fatalf("Expected %s, got %s", s, val)
		}
	}
}

func TestLongStrPerms(t *testing.T) {
	p := NewStringPermuter("a")

	expectNextStr(t, p, "")
	expectNextStr(t, p, "a")
	expectNextStr(t, p, "aa")
	expectNextStr(t, p, "aaa")
	expectNextStr(t, p, "aaaa")
	expectNextStr(t, p, "aaaaa")
	expectNextStr(t, p, "aaaaaa")
	expectNextStr(t, p, "aaaaaaa")
	expectNextStr(t, p, "aaaaaaaa")
	expectNextStr(t, p, "aaaaaaaaa")
	expectNextStr(t, p, "aaaaaaaaaa")
	expectNextStr(t, p, "aaaaaaaaaaa")
}
