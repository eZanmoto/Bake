// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package template

import (
	"testing"
)

func TestExpandVar(t *testing.T) {
	// Arrange
	before, after := "a", "b"
	d := &Dict{before: after}

	// Act & Assert
	testExpand(t, d, before, before)
	testExpand(t, d, "{"+before+"}", after)
	testExpand(t, d, before+"{"+before+"}", before+after)
	testExpand(t, d, "{"+before+"}"+before, after+before)
	testExpand(t, d, before+"{"+before+"}"+before, before+after+before)
}

func testExpand(t *testing.T, d *Dict, before, after string) {
	result, err := d.ExpandStr(before)
	if err != nil {
		t.Fatalf("Unexpected error while expanding '%s': %v",
			before, err)
	} else if after != result {
		t.Fatalf("Expanding:\n%s\nExpected:\n%s\nGot:\n%s\n",
			before, after, result)
	}
}
