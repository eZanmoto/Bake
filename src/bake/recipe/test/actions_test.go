// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestCommandBeforeBake(t *testing.T) {
	// Arrange
	command := parseAction(t, "descr\n touch wood")

	// don't want to dirty the build directory
	if err := cdNewTmp(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exists, err := fileExists("wood"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if exists {
		t.Fatalf("file 'wood' shouldn't exist before test")
	}

	// Act
	result, err := command.beforeBake()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.success() {
		t.Fatalf("unexpected failure: %s", result.detail())
	}

	if exists, err := fileExists("wood"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if !exists {
		t.Fatalf("file 'wood' should exist after test")
	}
}

func parseAction(t *testing.T, def string) testAction {
	buf := bytes.NewBufferString(def)

	tests, err := readTypeTests(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tests) != 1 {
		t.Fatalf("expected 1 parsed tests, got %d", len(tests))
	}

	actions := tests[0].actions()
	if len(actions) != 1 {
		t.Fatalf("expected 1 test action, got %d", len(actions))
	}

	return actions[0]
}

func cdNewTmp() error {
	tmpDirPath, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}

	return os.Chdir(tmpDirPath)
}

func fileExists(path string) (exists bool, err error) {
	_, err = os.Stat(path)

	exists = err == nil
	if os.IsNotExist(err) {
		err = nil
	}

	return
}

func TestCommandAfterBake(t *testing.T) {
	// Arrange
	command := parseAction(t, "descr\n touch wood")

	// don't want to dirty the build directory
	if err := cdNewTmp(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exists, err := fileExists("wood"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if exists {
		t.Fatalf("file 'wood' shouldn't exist before test")
	}

	// Act
	result, err := command.afterBake()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.success() {
		t.Fatalf("unexpected failure: %s", result.detail())
	}

	if exists, err := fileExists("wood"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if !exists {
		t.Fatalf("file 'wood' should exist after test")
	}
}
