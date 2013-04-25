// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package test provides methods for testing bake recipes.
package test

// In the test package, "test" is always a noun, never a verb. Instead of
// "testing", one "runs" a test.

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

const (
	testDirName = "tests"
)

// A typeTestGroup is a collection of tests for a set of bake project types.
//
// A bake project type is a type that is passed to bake to specify the type of
// project to be generated. Each instance of typeTestGroup contains tests for
// validating the behaviour of a single set of these project types.
type typeTestGroup struct {
	types []string
	tests []*typeTest
}

func (g *typeTestGroup) Types() []string {
	return g.types
}

func (g *typeTestGroup) Tests() []*typeTest {
	return g.tests
}

// Tests a recipe and returns true if the test succeeded.
func TestRecipe(lang string, recpDirPath string) (bool, error) {
	testDirPath := path.Join(recpDirPath, testDirName)

	typeTestScripts, err := ioutil.ReadDir(testDirPath)
	if err != nil {
		return false, err
	}

	groups := make([]*typeTestGroup, 0, len(typeTestScripts))
	for _, typeTestScript := range typeTestScripts {
		if typeTestScript.IsDir() {
			fmt.Printf("unexpected dir '%s' in '%s', skipping...",
				typeTestScript.Name(), testDirPath)
			continue
		}

		typeTestScriptName := typeTestScript.Name()

		typeTestScriptPath := path.Join(testDirPath, typeTestScriptName)

		typeTests, err := readTypeTestScript(typeTestScriptPath)
		if err != nil {
			return false, err
		}

		groups = append(groups, &typeTestGroup{
			strings.Split(typeTestScriptName, "_"),
			typeTests,
		})
	}

	var tempDir string
	if tempDir, err = ioutil.TempDir("", lang); err != nil {
		return false, err
	}

	passed := true
	for _, group := range groups {
		passed_, err := runTypeTestGroup(lang, tempDir, group)
		if err != nil {
			return false, err
		}
		passed = passed && passed_
	}

	return passed, nil
}
