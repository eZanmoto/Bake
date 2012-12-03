// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

import (
	"strings"
	"testing"
)

type inclTest struct {
	source   string
	expected *fsNode
}

var (
	inclTests = []inclTest{
		{"" +
			"a",
			&fsNode{"", []*fsNode{
				{"a", nil},
			}}},

		{"" +
			"a/\n" +
			"\tb",
			&fsNode{"", []*fsNode{
				{"a", []*fsNode{
					{"b", nil},
				}},
			}}},

		{"" +
			"a/\n" +
			"\tb\n" +
			"\tc",
			&fsNode{"", []*fsNode{
				{"a", []*fsNode{
					{"b", nil},
					{"c", nil},
				}},
			}}},

		{"" +
			"a/\n" +
			"\tb/\n" +
			"\t\tc",
			&fsNode{"", []*fsNode{
				{"a", []*fsNode{
					{"b", []*fsNode{
						{"c", nil},
					}},
				}},
			}}},

		{"" +
			"a/\n" +
			"\tb/\n" +
			"\t\tc\n" +
			"\td",
			&fsNode{"", []*fsNode{
				{"a", []*fsNode{
					{"b", []*fsNode{
						{"c", nil},
					}},
					{"d", nil},
				}},
			}}},
	}
)

func TestReadIncls(t *testing.T) {
	for _, test := range inclTests {
		result, err := parseIncls(strings.NewReader(test.source))
		if err != nil {
			t.Errorf("Failed: %v", err)
		} else if !result.equals(test.expected) {
			t.Errorf("\nParsed:\n" + test.source +
				"\nExpected:\n" + test.expected.String() +
				"\nGot:\n" + result.String())
		}
	}
}

func (n *fsNode) equals(m *fsNode) bool {
	if n.name != m.name {
		return false
	}

	if n.isDir() != m.isDir() {
		return false
	}

	if !n.isDir() {
		return true
	}

	for _, nChild := range n.children {
		mChild, ok := m.childNamed(nChild.name)

		if !ok || !nChild.equals(mChild) {
			return false
		}
	}

	return true
}
