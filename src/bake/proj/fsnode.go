// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

// fsnode.go contains the definition of the fsNode type, which represents a
// filesystem file or directory.

import (
	"sort"
	"strings"
)

type fsNode struct {
	name     string
	children []*fsNode
}

func (n *fsNode) isDir() bool {
	return n.children != nil
}

func (n *fsNode) childNamed(name string) (*fsNode, bool) {
	if !n.isDir() {
		return nil, false
	}

	for _, child := range n.children {
		if child.name == name {
			return child, true
		}
	}

	return nil, false
}

func (n *fsNode) addDir(name string) {
	if _, exists := n.childNamed(name); !exists {
		n.children = append(n.children, &fsNode{name, []*fsNode{}})
	}
}

func (n *fsNode) addFile(name string) {
	if _, exists := n.childNamed(name); !exists {
		n.children = append(n.children, &fsNode{name, nil})
	}
}

func (n *fsNode) String() string {
	s := n.name
	if n.children != nil {
		s += "/"

		names := make([]string, len(n.children))
		for i, c := range n.children {
			names[i] = c.name
		}

		sort.Strings(names)
		for _, name := range names {
			c, _ := n.childNamed(name)
			s += strings.Replace("\n"+c.String(), "\n", "\n\t", -1)
		}
	}
	return s
}
