// Copyright 2014 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package fs

import (
	"sort"
	"strings"
)

const (
	dirSep = "/"
)

type Node struct {
	isDir    bool
	name     string
	children []*Node
}

func NewFile(name string) *Node {
	return &Node{false, name, nil}
}

func NewDir(name string, children ...*Node) *Node {
	return &Node{true, name, children}
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) Children() []*Node {
	return n.children
}

func (n *Node) IsDir() bool {
	return n.isDir
}

func (n *Node) ChildNamed(name string) (*Node, bool) {
	if !n.IsDir() {
		return nil, false
	}

	for _, child := range n.children {
		if child.name == name {
			return child, true
		}
	}

	return nil, false
}

func (n *Node) AddNode(node *Node) *Node {
	if _, exists := n.ChildNamed(node.Name()); !exists {
		n.children = append(n.children, node)
	}
	return n
}

func (n *Node) AddDir(name string) *Node {
	if _, exists := n.ChildNamed(name); !exists {
		n.children = append(n.children, NewDir(name))
	}
	return n
}

func (n *Node) AddFile(name string) *Node {
	if _, exists := n.ChildNamed(name); !exists {
		n.children = append(n.children, NewFile(name))
	}
	return n
}

func (n *Node) String() string {
	s := n.name
	if n.children != nil {
		s += dirSep

		names := make([]string, len(n.children))
		for i, c := range n.children {
			names[i] = c.name
		}

		sort.Strings(names)
		for _, name := range names {
			c, _ := n.ChildNamed(name)
			s += strings.Replace("\n"+c.String(), "\n", "\n\t", -1)
		}
	}
	return s
}
