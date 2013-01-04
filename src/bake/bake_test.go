// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"testing"
	"tests/perm"
)

var (
	bakeProg       = path.Join(os.Getenv("BAKE"), "bin", "bake")
	supportedLangs = sort.StringSlice([]string{
		"go",
	})
	unknownLang = "unknown"
)

func TestMissingOwnerArg(t *testing.T) {
	cmd, _, errput := runBake(t, "-n", "x", "-l", "go")

	if cmd.ProcessState.Success() {
		t.Fatalf("bake exited successfully, expected failure")
	}

	if len(errput) == 0 {
		t.Fatalf("Expected error, stderr was empty")
	}
}

func runBake(t *testing.T, args ...string) (cmd *exec.Cmd, o, e string) {
	cmd = exec.Command(bakeProg, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Couldn't get bake stdout: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("Couldn't get bake stderr: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		t.Fatalf("Couldn't start bake: %v", err)
	}

	if o, err = readLines(stdout); err != nil {
		t.Fatalf("Error reading stdout: %v", err)
	}

	if e, err = readLines(stderr); err != nil {
		t.Fatalf("Error reading stderr: %v", err)
	}

	cmd.Wait()

	return
}

func TestMissingNameArg(t *testing.T) {
	cmd, _, errput := runBake(t, "-o", "x", "-l", "go")

	if cmd.ProcessState.Success() {
		t.Fatalf("bake exited successfully, expected failure")
	}

	if len(errput) == 0 {
		t.Fatalf("Expected error, stderr was empty")
	}
}

func TestMissingLanguageArg(t *testing.T) {
	cmd, _, errput := runBake(t, "-o", "x", "-n", "x")

	if cmd.ProcessState.Success() {
		t.Fatalf("bake exited successfully, expected failure")
	}

	if len(errput) == 0 {
		t.Fatalf("Expected error, stderr was empty")
	}
}

func TestUnknownLanguageArg(t *testing.T) {
	supportedLangs.Sort()
	if supportedLangs.Search(unknownLang) != len(supportedLangs) {
		t.Fatalf("%s is a supported language", unknownLang)
	}

	cmd, _, errput := runBake(t, "-o", "x", "-n", "x", "-l", unknownLang)

	if cmd.ProcessState.Success() {
		t.Fatalf("bake exited successfully, expected failure")
	}

	if len(errput) == 0 {
		t.Fatalf("Expected error, stderr was empty")
	}
}

func TestLanguagesArg(t *testing.T) {
	cmd, output, errput := runBake(t, "-L")

	if !cmd.ProcessState.Success() {
		t.Fatalf("bake did not exit successfully")
	}

	if len(errput) > 0 {
		t.Fatalf("Didn't expect error: %s", errput)
	}

	langs := sort.StringSlice(strings.Split(output, "\n"))
	langs = langs[:len(langs)-1] // trim extra newline
	if m, n := supportedLangs.Len(), langs.Len(); m != n {
		t.Fatalf("Expected %d supported languages, got %d: '%s'",
			m, n, strings.Join(langs, "','"))
	}

	langs.Sort()
	supportedLangs.Sort()
	for i := 0; i < supportedLangs.Len(); i++ {
		a := supportedLangs[i]
		b := langs[i]
		if !strings.EqualFold(a, b) {
			t.Fatalf("Expected support for '%s', got '%s'", a, b)
		}
	}
}

func readLines(in io.Reader) (string, error) {
	bufin := bufio.NewReader(in)
	lines := ""
	for {
		line, err := bufin.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		lines += line
	}
	return lines, nil
}

func TestGenProjDir(t *testing.T) {
	dir := os.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Error changing to directory: %v", err)
	}

	nameGen := perm.NewStringPermuter("abc")
	nameGen.Permute() // Skip empty name

	for _, lang := range supportedLangs {
		name := nameGen.Permute()

		if e := os.RemoveAll(name); e != nil && !os.IsNotExist(e) {
			t.Fatalf("Error removing '%s': %v", name, e)
		}

		bake(t, name, "owner", lang)

		if fi, err := os.Stat(name); err != nil {
			t.Fatalf("Error getting status of '%s': %v", name, err)
		} else if !fi.IsDir() {
			t.Errorf("Expected '%s' to be directory", name)
		}
	}
}

func bake(t *testing.T, name, owner, lang string, xs ...string) {
	args := append(xs, "-n", name, "-o", owner, "-l", lang)
	cmd, _, errput := runBake(t, args...)

	if len(errput) != 0 {
		t.Fatalf("stderr was not empty: %s", errput)
	}

	if !cmd.ProcessState.Success() {
		t.Fatalf("bake exited with error")
	}
}

func TestGenProjType(t *testing.T) {
	dir := os.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Error changing to directory: %v", err)
	}

	name := "example"

	// Test without -t bin
	if e := os.RemoveAll(name); e != nil && !os.IsNotExist(e) {
		t.Fatalf("Error removing '%s': %v", name, e)
	}

	bake(t, name, "owner", "go")

	fname := path.Join(name, "src", strings.ToLower(name)+".go")
	if _, err := os.Stat(fname); err == nil || !os.IsNotExist(err) {
		t.Fatalf("File '%s' should not exist: %v", fname, err)
	}

	// Test with -t bin
	if e := os.RemoveAll(name); e != nil && !os.IsNotExist(e) {
		t.Fatalf("Error removing '%s': %v", name, e)
	}

	bake(t, name, "owner", "go", "-t", "bin")

	fname = path.Join(name,
		"src",
		strings.ToLower(name),
		strings.ToLower(name)+".go")
	if fi, err := os.Stat(fname); err != nil {
		t.Fatalf("File '%s' should exist: %v", fname, err)
	} else if fi.IsDir() {
		t.Errorf("Expected '%s' to be a file", fname)
	}
}
