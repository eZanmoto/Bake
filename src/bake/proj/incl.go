// Copyright 2012-2014 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

import (
	"bufio"
	"fmt"
	"fs"
	"io"
	"os"
	"strings"
)

const (
	// A suffix that denotes a filesystem directory in include files, and is
	// independent of the actual directory separator used by the runtime
	// platform.
	inclDirSep = '/'
)

// Return a filesystem description composed of files described by each include
// file in `paths`.
func ParseInclFiles(paths ...string) (*fs.Node, error) {
	root := fs.NewDir("")

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

		if err = addIncl(root, in); err != nil {
			return nil, err
		}
	}

	return root, nil
}

// Add files described in the `reader` stream to the root node `n`.
func addIncl(n *fs.Node, reader io.Reader) error {
	in := bufio.NewReader(reader)
	nodePath := []*fs.Node{n}
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
			return fmt.Errorf("Empty name in %s/", curDir.Name())
		} else if !isValidFsName(name) {
			return fmt.Errorf("%s is not a valid name", name)
		} else if isDirName(name) {
			d := name[:len(name)-1]
			curDir.AddDir(d)
			dir, _ := curDir.ChildNamed(d)
			nodePath = append(nodePath, dir)
			enterDir = true
		} else {
			curDir.AddFile(name)
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

func indentLvl(s string) int {
	return len(s) - len(strings.TrimLeft(s, "\t"))
}
