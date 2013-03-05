// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package diff

import (
	"fmt"
)

// ChangeType represents how a line differs between two string lists.
type ChangeType byte

const (
	Same ChangeType = 0
	Rem  ChangeType = 1
	Add  ChangeType = 2
)

// ChangeList represents how two string lists differ.
type ChangeList struct {
	stats []ChangeType
	lines []string
}

// Add adds a line and how it differs from an original line.
func (c *ChangeList) Add(stat ChangeType, line string) {
	c.stats = append(c.stats, stat)
	c.lines = append(c.lines, line)
}

// Len returns the total number of lines in this ChangeList.
func (c *ChangeList) Len() int {
	return len(c.stats)
}

// Get returns the nth change of this ChangeList.
func (c *ChangeList) Get(n int) (stat ChangeType, line string, err error) {
	if n < 0 || n >= c.Len() {
		err = fmt.Errorf("Index (%d) out of range [0,%d]", n, c.Len())
	} else {
		stat = c.stats[n]
		line = c.lines[n]
	}

	return
}
