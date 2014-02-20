// Copyright 2012-2014 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

import (
	"bake/env"
	"bufio"
	"fmt"
	"fs"
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
	root := fs.NewDir("{ProjectName}", incls.Children()...)

	return p.genDirConts(fs.NewDir("").AddNode(root), langRoot, "")
}

func joinAll(dir string, fnames []string) []string {
	paths := make([]string, len(fnames))
	copy(paths, fnames)
	for i, p := range paths {
		paths[i] = path.Join(dir, p)
	}
	return paths
}

func (p *Project) genDirConts(dir *fs.Node, srcDir, tgtDir string) error {
	for _, node := range dir.Children() {
		src := path.Join(srcDir, node.Name())

		tgtName, err := p.dict.ExpandStr(node.Name())
		if err != nil {
			return err
		}
		tgt := path.Join(tgtDir, tgtName)

		if node.Children() == nil { // not a dir
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
		if p.verbose {
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

	err = p.dict.Expand(in, out)
	if err != nil {
		return fmt.Errorf("%s%v", src, err)
	}

	if p.verbose {
		fmt.Printf("Generated file '%s'\n", tgt)
	} else {
		fmt.Printf("%s\n", tgt)
	}

	return err
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
