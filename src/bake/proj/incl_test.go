// Copyright 2012-2014 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

import (
	"fs"
	"io"
	"strings"
	"testing"
)

type inclTest struct {
	source   string
	expected *fs.Node
}

var (
	inclTests = []inclTest{
		{"",
			fs.NewDir(""),
		},

		{"" +
			"a",
			fs.NewDir("",
				fs.NewFile("a"),
			),
		},

		{"" +
			"a\n",
			fs.NewDir("",
				fs.NewFile("a"),
			),
		},

		{"" +
			"a/\n" +
			"\tb",
			fs.NewDir("",
				fs.NewDir("a",
					fs.NewFile("b"),
				),
			),
		},

		{"" +
			"a/\n" +
			"\tb\n",
			fs.NewDir("",
				fs.NewDir("a",
					fs.NewFile("b"),
				),
			),
		},

		{"" +
			"a/\n" +
			"\tb\n" +
			"\tc",
			fs.NewDir("",
				fs.NewDir("a",
					fs.NewFile("b"),
					fs.NewFile("c"),
				),
			),
		},

		{"" +
			"a/\n" +
			"\tb\n" +
			"\tc\n",
			fs.NewDir("",
				fs.NewDir("a",
					fs.NewFile("b"),
					fs.NewFile("c"),
				),
			),
		},

		{"" +
			"a/\n" +
			"\tb/\n" +
			"\t\tc",
			fs.NewDir("",
				fs.NewDir("a",
					fs.NewDir("b",
						fs.NewFile("c"),
					),
				),
			),
		},

		{"" +
			"a/\n" +
			"\tb/\n" +
			"\t\tc\n",
			fs.NewDir("",
				fs.NewDir("a",
					fs.NewDir("b",
						fs.NewFile("c"),
					),
				),
			),
		},

		{"" +
			"a/\n" +
			"\tb/\n" +
			"\t\tc\n" +
			"\td",
			fs.NewDir("",
				fs.NewDir("a",
					fs.NewDir("b",
						fs.NewFile("c"),
					),
					fs.NewFile("d"),
				),
			),
		},

		{"" +
			"a/\n" +
			"\tb/\n" +
			"\t\tc\n" +
			"\td\n",
			fs.NewDir("",
				fs.NewDir("a",
					fs.NewDir("b",
						fs.NewFile("c"),
					),
					fs.NewFile("d"),
				),
			),
		},
	}
)

func TestReadIncl(t *testing.T) {
	for _, test := range inclTests {
		result, err := parseIncl(strings.NewReader(test.source))
		if err != nil {
			t.Errorf("Failed: %v", err)
		} else if !equal(result, test.expected) {
			t.Errorf("\nParsed:\n%s\nExpected:\n%s\nGot:\n%s",
				test.source, test.expected.String(),
				result.String())
		}
	}
}

func parseIncl(reader io.Reader) (*fs.Node, error) {
	n := fs.NewDir("")
	if err := addIncl(n, reader); err != nil {
		return nil, err
	}
	return n, nil
}

func equal(n *fs.Node, m *fs.Node) bool {
	if n.Name() != m.Name() {
		return false
	}

	if n.IsDir() != m.IsDir() {
		return false
	}

	if !n.IsDir() {
		return true
	}

	for _, nChild := range n.Children() {
		mChild, ok := m.ChildNamed(nChild.Name())

		if !ok || !equal(nChild, mChild) {
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
	n := fs.NewDir("")
	sources := []string{
		"a/\n\tb/\n\t\tc\n",
		"a/\n\tb/\n\t\td\n",
		"a/\n\te/\n\t\tf\n",
	}
	expect := fs.NewDir("",
		fs.NewDir("a",
			fs.NewDir("b",
				fs.NewFile("c"),
				fs.NewFile("d"),
			),
			fs.NewDir("e",
				fs.NewFile("f"),
			),
		),
	)

	for _, source := range sources {
		if err := addIncl(n, strings.NewReader(source)); err != nil {
			t.Errorf("Failed: %v", err)
		}
	}

	if !equal(n, expect) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expect.String(), n.String())
	}
}
