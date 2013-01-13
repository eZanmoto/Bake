// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package perm provides utilities for creating permutations.
package perm

type Permuter interface {
	// Permute returns the next permutation of a sequence
	Permute() []byte // not interface{}, since it won't work with primitives
}
