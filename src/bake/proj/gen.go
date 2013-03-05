// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

import (
	"bake/env"
	"bufio"
	"diff"
	"fmt"
	"io"
	"os"
	"path"
)

const (
	baseInclFile = "base"
)

// GenTo generates the project p to dest.
func (p *Project) GenTo(dest string) error {
	templPath, err := env.TemplatesPath()
	if err != nil {
		return err
	}

	langRoot := path.Join(templPath, p.lang)

	filePaths := joinAll(langRoot, append(p.types, baseInclFile))
	incls, err := ParseInclFiles(filePaths...)
	if err != nil {
		return err
	}
	incls.name = "{ProjectName}"

	return p.genDirConts(&fsNode{children: []*fsNode{incls}}, langRoot, "")
}

func joinAll(dir string, fnames []string) []string {
	paths := make([]string, len(fnames))
	copy(paths, fnames)
	for i, p := range paths {
		paths[i] = path.Join(dir, p)
	}
	return paths
}

func (p *Project) genDirConts(dir *fsNode, srcDir, tgtDir string) error {
	for _, node := range dir.children {
		src := path.Join(srcDir, node.name)

		tgtName, err := p.parseStr(node.name)
		if err != nil {
			return err
		}
		tgt := path.Join(tgtDir, tgtName)

		if node.children == nil { // not a dir
			err = p.genFile(src, tgt)
		} else if err = p.genDir(tgt); err == nil {
			err = p.genDirConts(node, src, tgt)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Project) genFile(src, tgt string) error {
	out, err := os.OpenFile(tgt, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
		if p.resolve {
			if err := p.addUpdates(src, tgt); err != nil {
				return err
			}
			if p.verbose {
				defer fmt.Printf("Merged file '%s'\n", tgt)
			} else {
				defer fmt.Printf("%s\n", tgt)
			}
		} else if p.verbose {
			fmt.Printf("File '%s' exists, skipping...\n", tgt)
		}
		return nil
	}
	defer out.Close()

	in, err := os.OpenFile(src, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer in.Close()

	if p.verbose {
		defer fmt.Printf("Generated file '%s'\n", tgt)
	} else {
		defer fmt.Printf("%s\n", tgt)
	}

	return p.parse(in, out)
}

func (p *Project) genDir(dir string) error {
	if err := os.Mkdir(dir, 0777); err != nil {
		if !os.IsExist(err) {
			return err
		}
		if p.verbose {
			fmt.Printf("Directory '%s/' exists, skipping...\n", dir)
		}
	} else if p.verbose {
		fmt.Printf("Created directory '%s/'\n", dir)
	} else {
		fmt.Printf("%s/\n", dir)
	}

	return nil
}

func (p *Project) addUpdates(src, tgt string) error {
	tgtFile, err := os.OpenFile(tgt, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer tgtFile.Close()

	srcFile, err := os.OpenFile(src, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	chs := diff.Diff(readLines(srcFile), readLines(tgtFile))

	lastLineWasNew := false
	out := bufio.NewWriter(tgtFile)

	for i := 0; i < chs.Len(); i++ {
		stat, line, err := chs.Get(i)
		if err != nil {
			return err
		}

		lineIsNew := stat == diff.Add
		if lastLineWasNew != lineIsNew {
			lastLineWasNew = lineIsNew

			if lineIsNew {
				out.WriteString("++++ <new>\n")
			} else {
				out.WriteString("++++ </new>\n")
			}
		}

		out.WriteString(line + "\n")
	}
	out.Flush()

	return nil
}

func readLines(reader io.Reader) []string {
	var lines []string
	var err error
	in := bufio.NewReader(reader)

	for err != io.EOF {
		isPrefix := true
		line := ""
		for isPrefix && err != io.EOF {
			var bytes []byte
			bytes, isPrefix, err = in.ReadLine()
			line += string(bytes)
		}
		lines = append(lines, line)
	}

	return lines
}
