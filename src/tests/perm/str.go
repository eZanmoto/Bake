// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package perm

// A primitive implementation of a string permuter
type StringPermuter struct {
	permer *BytePermuter
	chars  []byte
}

func NewStringPermuter(s string) *StringPermuter {
	chars := []byte(s)
	p := NewBytePermuter(0, byte(len(chars)))
	return &StringPermuter{p, chars}
}

func (p *StringPermuter) Permute() string {
	cs := p.permer.Permute()
	for i, c := range cs {
		cs[i] = p.chars[c]
	}
	return string(cs)
}
