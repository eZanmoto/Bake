// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

import (
	"bake/env"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// GenTo generates the project p to dest.
func (p *Project) GenTo(dest string) error {
	templPath, err := env.TemplatesPath()

	if err != nil {
		return err
	}

	root := path.Join(templPath, p.lang)
	return p.genDir(root, dest, "{ProjectName}")
}

func (p *Project) genDir(srcDir, tgtDir, dirName string) error {
	src := path.Join(srcDir, dirName)

	tgtName, err := p.parseStr(dirName)
	if err != nil {
		return err
	}
	tgt := path.Join(tgtDir, tgtName)

	if err := os.Mkdir(tgt, 0777); err != nil {
		if !os.IsExist(err) {
			return err
		}
		if p.verbose {
			fmt.Printf("Directory '%s' exists, skipping...\n", tgt)
		}
	} else if p.verbose {
		fmt.Printf("Created directory '%s'\n", tgt)
	} else {
		fmt.Printf("%s\n", tgt)
	}

	return p.genDirConts(src, tgt)
}

func (p *Project) genDirConts(srcDir, tgtDir string) error {
	files, err := ioutil.ReadDir(srcDir)

	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			err = p.genDir(srcDir, tgtDir, file.Name())
		} else {
			err = p.genFile(srcDir, tgtDir, file.Name())
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Project) genFile(srcDir, tgtDir, fileName string) error {
	src := path.Join(srcDir, fileName)

	tgtName, err := p.parseStr(fileName)
	if err != nil {
		return err
	}
	tgt := path.Join(tgtDir, tgtName)

	in, err := os.OpenFile(src, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer in.Close()

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

	defer fmt.Printf("Generated file '%s'\n", tgt)

	return p.parse(in, out)
}
