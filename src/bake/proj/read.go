// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

import (
	"bufio"
	"fmt"
	"io"
	"os"
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

func (n *fsNode) addDir(name string) bool {
	if _, exists := n.childNamed(name); !exists {
		n.children = append(n.children,
			&fsNode{name, []*fsNode{}})
		return true
	}
	return false
}

func (n *fsNode) addFile(name string) bool {
	if _, exists := n.childNamed(name); !exists {
		n.children = append(n.children,
			&fsNode{name, nil})
		return true
	}
	return false
}

func (n *fsNode) String() string {
	str := n.name + "/"
	if n.children != nil {
		for _, child := range n.children {
			str += strings.Replace("\n"+child.String(),
				"\n", "\n\t", -1)
		}
	}
	return str
}

func parseInclsFile(path string) (*fsNode, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	in := bufio.NewReader(file)
	_, err = in.ReadString('\n') // Skip description
	if err != nil {
		return nil, err
	}

	return parseIncls(in)
}

func parseIncls(reader io.Reader) (*fsNode, error) {
	in := bufio.NewReader(reader)
	nodePath := []*fsNode{{"", []*fsNode{}}}

	for {
		var line string
		line, err := in.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		lvl := indentLvl(line)
		if lvl >= len(nodePath) {
			return nil, fmt.Errorf("Bad indentation: '%s'", line)
		} else {
			nodePath = nodePath[:lvl+1]
		}
		curDir := nodePath[len(nodePath)-1]

		name := strings.TrimRight(line, "\n")[lvl:]
		if len(name) == 0 {
			return nil, fmt.Errorf("Empty name in %s/", curDir.name)
		} else if name[len(name)-1] == inclDirSep {
			d := name[:len(name)-1]
			if !curDir.addDir(d) {
				return nil, fmt.Errorf("Repeated name '%s'", d)
			}
			dir, _ := curDir.childNamed(d)
			nodePath = append(nodePath, dir)
		} else if !curDir.addFile(name) {
			return nil, fmt.Errorf("Repeated name '%s'", name)
		}
	}

	return nodePath[0], nil
}

func indentLvl(s string) int {
	return len(s) - len(strings.TrimLeft(s, "\t"))
}
