// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package diff provides utilities for finding differences in batches of lines.
package diff

import (
	"strings"
)

// Diff returns a list of changes that can transform a into b.
func Diff(a, b []string) ChangeList {
	diff := &ChangeList{}
	aIndex := 0

	for _, bLine := range b {
		if aIndex >= len(a) {
			diff.Add(Add, bLine)
		} else if a[aIndex] == bLine {
			diff.Add(Same, a[aIndex])
			aIndex++
		} else if len(strings.TrimSpace(bLine)) == 0 {
			diff.Add(Add, bLine)
		} else {
			if i := indexOf(a[aIndex:], bLine); i == -1 {
				diff.Add(Add, bLine)
			} else {
				for j := 0; j < i; j++ {
					diff.Add(Rem, a[aIndex+j])
				}
				aIndex += i
				diff.Add(Same, a[aIndex])
				aIndex += 1
			}
		}
	}

	for aIndex < len(a) {
		diff.Add(Rem, a[aIndex])
		aIndex++
	}

	return *diff
}

func indexOf(xs []string, tgt string) int {
	for i, x := range xs {
		if x == tgt {
			return i
		}
	}
	return -1
}
