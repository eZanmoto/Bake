// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

import (
	"io"
	"strings"
	"testing"
)

type inclTest struct {
	source   string
	expected *fsNode
}

var (
	inclTests = []inclTest{
		{"",
			&fsNode{"", []*fsNode{}},
		},

		{"" +
			"a",
			&fsNode{"", []*fsNode{
				{"a", nil},
			}},
		},

		{"" +
			"a\n",
			&fsNode{"", []*fsNode{
				{"a", nil},
			}},
		},

		{"" +
			"a/\n" +
			"\tb",
			&fsNode{"", []*fsNode{
				{"a", []*fsNode{
					{"b", nil},
				}},
			}},
		},

		{"" +
			"a/\n" +
			"\tb\n",
			&fsNode{"", []*fsNode{
				{"a", []*fsNode{
					{"b", nil},
				}},
			}},
		},

		{"" +
			"a/\n" +
			"\tb\n" +
			"\tc",
			&fsNode{"", []*fsNode{
				{"a", []*fsNode{
					{"b", nil},
					{"c", nil},
				}},
			}},
		},

		{"" +
			"a/\n" +
			"\tb\n" +
			"\tc\n",
			&fsNode{"", []*fsNode{
				{"a", []*fsNode{
					{"b", nil},
					{"c", nil},
				}},
			}},
		},

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
			}},
		},

		{"" +
			"a/\n" +
			"\tb/\n" +
			"\t\tc\n",
			&fsNode{"", []*fsNode{
				{"a", []*fsNode{
					{"b", []*fsNode{
						{"c", nil},
					}},
				}},
			}},
		},

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
			}},
		},

		{"" +
			"a/\n" +
			"\tb/\n" +
			"\t\tc\n" +
			"\td\n",
			&fsNode{"", []*fsNode{
				{"a", []*fsNode{
					{"b", []*fsNode{
						{"c", nil},
					}},
					{"d", nil},
				}},
			}},
		},
	}
)

func TestReadIncl(t *testing.T) {
	for _, test := range inclTests {
		result, err := parseIncl(strings.NewReader(test.source))
		if err != nil {
			t.Errorf("Failed: %v", err)
		} else if !result.equals(test.expected) {
			t.Errorf("\nParsed:\n%s\nExpected:\n%s\nGot:\n%s",
				test.source, test.expected.String(),
				result.String())
		}
	}
}

func parseIncl(reader io.Reader) (*fsNode, error) {
	n := newRootDir()
	if err := n.addIncl(reader); err != nil {
		return nil, err
	}
	return n, nil
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

func TestEmptyDirFail(t *testing.T) {
	expectFail(t, ""+
		"a/\n")

	expectFail(t, ""+
		"a/")

	expectFail(t, ""+
		"a/\n"+
		"b\n")

	expectFail(t, ""+
		"a/\n"+
		"b")
}

func expectFail(t *testing.T, src string) {
	if result, err := parseIncl(strings.NewReader(src)); err == nil {
		t.Errorf("Expected failure parsing:\n%s\nGot:\n%s", src, result)
	}
}

func TestBadIndentation(t *testing.T) {
	expectFail(t, ""+
		"a/\n"+
		"\t\tb\n")

	expectFail(t, ""+
		"a/\n"+
		"\t\tb")
}

func TestAddIncl(t *testing.T) {
	n := newRootDir()
	sources := []string{
		"a/\n\tb/\n\t\tc\n",
		"a/\n\tb/\n\t\td\n",
		"a/\n\te/\n\t\tf\n",
	}
	expect := &fsNode{"", []*fsNode{
		{"a", []*fsNode{
			{"b", []*fsNode{
				{"c", nil},
				{"d", nil},
			}},
			{"e", []*fsNode{
				{"f", nil},
			}},
		}},
	}}

	for _, source := range sources {
		if err := n.addIncl(strings.NewReader(source)); err != nil {
			t.Errorf("Failed: %v", err)
		}
	}

	if !expect.equals(n) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expect.String(), n.String())
	}
}
