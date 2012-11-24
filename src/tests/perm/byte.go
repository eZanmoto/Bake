// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package perm

const (
	defaultSliceSize = 8
)

// A primitive implementation of a byte permuter
type BytePermuter struct {
	bytes []byte
	start byte
	count byte
}

// NewBytePermuter creates a structure that creates byte slice permutations.
// NewBytePermuter will panic on count of 0, as one can't permute over a range
// of 0 different numbers.
func NewBytePermuter(start, count byte) *BytePermuter {
	if count == 0 {
		panic("Can't permute over 0 different values")
	}
	return &BytePermuter{make([]byte, 0, defaultSliceSize), start, count}
}

func (p *BytePermuter) Permute() []byte {
	bs := p.bytes
	last := make([]byte, len(bs))

	n := copy(last, bs)
	if n != len(bs) {
		panic("Error copying bytes")
	}

	overflow := true

	for i := 0; i < len(bs) && overflow; i++ {
		if bs[i]-p.start == p.count-1 {
			bs[i] = p.start
		} else {
			overflow = false
			bs[i]++
		}
	}

	if overflow {
		n := len(bs)
		if n == cap(bs) {
			cs := make([]byte, n*2)
			for i := range bs {
				cs[i] = bs[i]
			}
			bs = cs
		}
		bs = bs[0 : n+1]
		bs[n] = p.start
	}

	p.bytes = bs

	return last
}
