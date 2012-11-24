// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package perm

import "testing"

func TestFirstPerms(t *testing.T) {
	p := NewBytePermuter(0, 3)

	expectNext(t, p)
	expectNext(t, p, 0)
	expectNext(t, p, 1)
	expectNext(t, p, 2)
	expectNext(t, p, 0, 0)
	expectNext(t, p, 1, 0)
	expectNext(t, p, 2, 0)
	expectNext(t, p, 0, 1)
	expectNext(t, p, 1, 1)
	expectNext(t, p, 2, 1)
	expectNext(t, p, 0, 2)
	expectNext(t, p, 1, 2)
	expectNext(t, p, 2, 2)
	expectNext(t, p, 0, 0, 0)
}

func expectNext(t *testing.T, p Permuter, b ...byte) {
	val := p.Permute()

	if len(val) != len(b) {
		t.Fatalf("Expected %d bytes, got %d", len(b), len(val))
	}

	for i := 0; i < len(val); i++ {
		if val[i] != b[i] {
			t.Fatalf("Expected %d, got %d", b[i], val[i])
		}
	}
}

func TestLongPerms(t *testing.T) {
	p := NewBytePermuter(0, 1)

	expectNext(t, p)
	expectNext(t, p, 0)
	expectNext(t, p, 0, 0)
	expectNext(t, p, 0, 0, 0)
	expectNext(t, p, 0, 0, 0, 0)
	expectNext(t, p, 0, 0, 0, 0, 0)
	expectNext(t, p, 0, 0, 0, 0, 0, 0)
	expectNext(t, p, 0, 0, 0, 0, 0, 0, 0)
	expectNext(t, p, 0, 0, 0, 0, 0, 0, 0, 0)
	expectNext(t, p, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	expectNext(t, p, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	expectNext(t, p, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
}

func TestPermCount(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Expected error, got none")
		}
	}()

	NewBytePermuter(0, 0)
}
