// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

const (
	inclDirSep = '/' // This is independent of the actual platform
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

type NodePath []*fsNode

func ParseInclFiles(paths ...string) (*fsNode, error) {
	root := newRootDir()

	for _, path := range paths {
		file, err := os.OpenFile(path, os.O_RDONLY, 0666)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		in := bufio.NewReader(file)

		// Skip description
		if _, err = in.ReadString('\n'); err != nil {
			return nil, err
		}

		if err = root.addIncl(in); err != nil {
			return nil, err
		}
	}

	return root, nil
}

func (n *fsNode) addIncl(reader io.Reader) error {
	in := bufio.NewReader(reader)
	nodePath := []*fsNode{n}
	enterDir := false

	for {
		line, err := in.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			if len(line) == 0 {
				break
			}
		}

		lvl := indentLvl(line)
		if enterDir && lvl != len(nodePath)-1 || lvl >= len(nodePath) {
			return fmt.Errorf("Bad indentation: '%s'", line)
		} else {
			nodePath = nodePath[:lvl+1]
		}
		curDir := nodePath[len(nodePath)-1]
		enterDir = false

		name := strings.TrimRight(line, "\n\r")[lvl:]
		if len(name) == 0 {
			return fmt.Errorf("Empty name in %s/", curDir.name)
		} else if !isValidFsName(name) {
			return fmt.Errorf("%s is not a valid name", name)
		} else if isDirName(name) {
			d := name[:len(name)-1]
			curDir.addDir(d)
			dir, _ := curDir.childNamed(d)
			nodePath = append(nodePath, dir)
			enterDir = true
		} else {
			curDir.addFile(name)
		}

		if err == io.EOF {
			break
		}
	}

	if enterDir {
		return fmt.Errorf("Expected dir contents, got EOF")
	}

	return nil
}

func isDirName(d string) bool {
	return d[len(d)-1] == inclDirSep
}

func isValidFsName(n string) bool {
	for i := 0; i < len(n)-1; i++ {
		if n[i] == inclDirSep {
			return false
		}
	}
	return true
}

func newRootDir() *fsNode {
	return &fsNode{"", []*fsNode{}}
}

func indentLvl(s string) int {
	return len(s) - len(strings.TrimLeft(s, "\t"))
}
