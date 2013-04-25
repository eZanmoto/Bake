// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package test

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"strio"
)

const (
	testDirPerm = 0777
)

const (
	testDirectiveIndex = 0
)

const (
	beforeBakeDir = "pre"
	bakeTestDir   = "test"
)

const (
	logFileName = "log"
)

func readTypeTestScript(scriptPath string) ([]*typeTest, error) {
	script, err := os.Open(scriptPath)
	if err != nil {
		return nil, err
	}
	defer script.Close()

	in := newLineCountReader(strio.NewLineReader(script))

	tests, err := readTypeTests(in)
	if err != nil {
		err = fmt.Errorf("%s:%d:\t%v", scriptPath, in.LineNum()-1, err)
	}

	for _, test := range tests {
		test.loc_ = scriptPath + ":" + test.loc_
	}

	return tests, err
}

type lineCountReader struct {
	in      strio.LineReader
	lineNum int
}

func newLineCountReader(in strio.LineReader) *lineCountReader {
	return &lineCountReader{in, 1}
}

func (r *lineCountReader) ChompLine() (string, error) {
	r.lineNum++
	s, e := r.in.ChompLine()
	return s, e
}

func (r *lineCountReader) ReadLine() (string, error) {
	r.lineNum++
	return r.in.ReadLine()
}

func (r *lineCountReader) LineNum() int {
	return r.lineNum
}

func readTypeTests(reader strio.LineReader) ([]*typeTest, error) {
	var err error

	in := newLineCountReader(reader)

	tests := make([]*typeTest, 0, 1)
	for {
		lineNum := strconv.Itoa(in.LineNum())

		var descr string
		descr, err = in.ChompLine()
		if err != nil {
			if err == io.EOF {
				err = errors.New("file must end with newline")
			}
			break
		}

		if len(descr) == 0 {
			err = errors.New("description cannot be empty")
			break
		}

		var actions []testAction
		actions, err = readTypeTestCommands(in)
		tests = append(tests, &typeTest{lineNum, descr, actions})

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
	}
	return tests, err
}

type typeTest struct {
	loc_     string       // A description of where the test is located
	descr_   string       // A description of what the test asserts
	actions_ []testAction // The actions that the test performs
}

func (t *typeTest) loc() string {
	return t.loc_
}

func (t *typeTest) descr() string {
	return t.descr_
}

func (t *typeTest) actions() []testAction {
	return t.actions_
}

func readTypeTestCommands(in strio.LineReader) ([]testAction, error) {
	var err error
	actions := make([]testAction, 0, 1)

	for {
		var line string
		line, err = in.ChompLine()

		if len(line) == 0 {
			break
		}

		if err != nil {
			if err == io.EOF {
				err = fmt.Errorf("file must end with newline")
			}
			break
		}
		// inv: len(line) > 0 && (err == nil || err == io.EOF)

		cmd := line[testDirectiveIndex+1:]

		action := parseTestAction(rune(line[testDirectiveIndex]), cmd)
		if action == nil {
			err = fmt.Errorf("'%c' is not a valid test directive",
				line[testDirectiveIndex])
			break
		}
		actions = append(actions, action)
	}

	return actions, err
}

func parseTestAction(actionSpecifier rune, cmd string) testAction {
	var action testAction

	// this is the only place the test directive constants are used, so they
	// are defined using magic constants since it is easier to link the
	// symbols to the actions they should perform
	switch actionSpecifier {
	case ' ':
		action = newCommand(cmd)
	case '+':
		action = newPass(cmd)
	case '=':
		action = newBuildPass(cmd)
	default:
		action = nil
	}

	return action
}

func runTypeTestGroup(lang string, testDirPath string,
	group *typeTestGroup) (passed bool, err error) {

	passed = true

	typeTestDirName := strings.Join(group.Types(), "_")
	typeTestDirPath := path.Join(testDirPath, typeTestDirName)
	if err = os.Mkdir(typeTestDirPath, testDirPerm); err != nil {
		return
	}

	projName := "Project"
	for _, test := range group.Tests() {
		for _, action := range test.actions() {
			action.AddVars(map[string]string{
				"ProjectName":      projName,
				"ProjectNameLower": strings.ToLower(projName),
			})
		}
	}

	testDir := path.Join(typeTestDirPath, logFileName)
	var out *os.File
	out, err = os.Create(testDir)
	for _, test := range group.Tests() {
		for _, action := range test.actions() {
			action.AddVars(map[string]string{"TestDir": testDir})
		}
	}
	if err != nil {
		err = fmt.Errorf("error opening log: %v", err)
		return
	}
	defer out.Close()
	log := &logr{out, os.Stderr}

	if err = cdNewDir(path.Join(typeTestDirPath, beforeBakeDir)); err != nil {
		return
	}
	for _, test := range group.Tests() {
		for _, action := range test.actions() {
			r, err := action.beforeBake()

			printResultToLog(test, r, err, log)

			if err != nil || !r.success() {
				passed = false
				break
			}
		}
	}
	log.Printf("\n")

	testDir = path.Join(typeTestDirPath, bakeTestDir)
	for _, test := range group.Tests() {
		for _, action := range test.actions() {
			action.AddVars(map[string]string{"TestDir": testDir})
		}
	}
	if err = cdNewDir(testDir); err != nil {
		return
	}
	if err = bakeWithLog(projName, lang, group.Types(), log); err != nil {
		return
	}
	for _, test := range group.Tests() {
		for _, action := range test.actions() {
			r, err := action.afterBake()

			printResultToLog(test, r, err, log)

			if err != nil || !r.success() {
				passed = false
				break
			}
		}
	}

	return
}

type Printfer interface {
	Printf(format string, v ...interface{})
}

func stdlogger(out io.Writer) Printfer {
	return log.New(out, "", 0)
}

type Errorfer interface {
	Errorf(format string, v ...interface{})
}

type Logger interface {
	Printfer
	Errorfer
}

// logr directs Printf output to the standard output stream and Errorf output to
// both the standard output stream and the error output stream.
type logr struct {
	out io.Writer
	err io.Writer
}

const (
	outstream = 1 << iota
	errstream
)

func (l *logr) Printf(format string, v ...interface{}) {
	l.dmux(outstream, format, v...)
}

func (l *logr) Errorf(format string, v ...interface{}) {
	l.dmux(outstream|errstream, format, v...)
}

func (l *logr) dmux(streams int, format string, v ...interface{}) {
	s := []byte(fmt.Sprintf(format, v...))

	if streams&outstream != 0 {
		l.out.Write(s)
	}

	if streams&errstream != 0 {
		l.err.Write(s)
	}
}

func cdNewDir(dir string) error {
	if err := os.Mkdir(dir, testDirPerm); err != nil {
		return err
	}
	return os.Chdir(dir)
}

func printResultToLog(t *typeTest, r *result, err error, out Logger) {
	if err != nil {
		out.Errorf("--- ERROR: %s\n", t.descr())
		out.Errorf("%s:\t%s\n", t.loc(), err)
	} else {
		stat := "PASS"
		outf := func(f string, v ...interface{}) { out.Printf(f, v...) }
		if !r.success() {
			stat = "FAIL"
			outf = func(f string, v ...interface{}) { out.Errorf(f, v...) }
		}
		outf("--- %s: %s\n", stat, t.descr())
		outf("%s:\t%s\n", t.loc(), r.descr())
		outf("\t%s\n", strings.Replace(r.detail(), "\n", "\n\t", -1))
	}
}

func bakeWithLog(name, lang string, types []string, out Printfer) error {
	typeString := "base"
	if len(types) != 0 {
		typeString = strings.Join(types, ",")
	}

	cmd := exec.Command(
		path.Join(os.Getenv("BAKE"), "bin", "bake"),
		"-v",
		"-o", "Owner",
		"-l", lang,
		"-n", name,
		"-t", typeString,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("couldn't get bake stdout: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("couldn't get bake stderr: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("couldn't start bake: %v", err)
	}

	output, err := strio.ReadAll(stdout)
	if err != nil {
		return fmt.Errorf("error reading bake stdout: %v", err)
	}
	out.Printf("%s\n", output)

	if errput, err := strio.ReadAll(stderr); err != nil {
		return fmt.Errorf("error reading bake stderr: %v", err)
	} else if len(errput) != 0 {
		return fmt.Errorf("unexpected error output: %s", errput)
	}

	return cmd.Wait()
}
