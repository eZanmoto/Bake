// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package test

import (
	"bytes"
	"strio"
	"testing"
)

const (
	commandDirective   = ' '
	passDirective      = '+'
	buildPassDirective = '='
)

const (
	goodCmd = iota
	badCmd
	errCmd
)

const (
	succeeds = iota
	fails
	errs
)

const (
	beforeBake = true
	afterBake  = false
)

// The following format should be used for the action tests:
//
// func TestCommandActionWithGoodCmdSucceedsBeforeBake(t *testing.T) {
// 	assert(t, commandDirective, goodCmd, succeeds, beforeBake)
// }
//
// pros: will know exactly what asserts failed
// cons: redundant naming scheme -- can solve by generating tests from asserts
//
// TODO generate tests

func TestCommandAction(t *testing.T) {
	assert(t, commandDirective, goodCmd, succeeds, beforeBake)
	assert(t, commandDirective, badCmd, fails, beforeBake)
	assert(t, commandDirective, errCmd, errs, beforeBake)

	assert(t, commandDirective, goodCmd, succeeds, afterBake)
	assert(t, commandDirective, badCmd, fails, afterBake)
	assert(t, commandDirective, errCmd, errs, afterBake)
}

func TestPassAction(t *testing.T) {
	assert(t, passDirective, badCmd, succeeds, beforeBake)
	assert(t, passDirective, goodCmd, fails, beforeBake)
	assert(t, passDirective, errCmd, errs, beforeBake)

	assert(t, passDirective, goodCmd, succeeds, afterBake)
	assert(t, passDirective, badCmd, fails, afterBake)
	assert(t, passDirective, errCmd, errs, afterBake)
}

func TestBuildPassAction(t *testing.T) {
	assert(t, buildPassDirective, goodCmd, fails, beforeBake)
	assert(t, buildPassDirective, badCmd, fails, beforeBake)
	assert(t, buildPassDirective, errCmd, succeeds, beforeBake)

	assert(t, buildPassDirective, goodCmd, succeeds, afterBake)
	assert(t, buildPassDirective, badCmd, fails, afterBake)
	assert(t, buildPassDirective, errCmd, errs, afterBake)
}

func assert(t *testing.T, directive rune, cmd uint, expect uint,
	beforeBake bool) {

	// While this method of testing is not the best approach, it will be
	// kept in place until the tests are generated.
	//
	// This method is so big because it's doing multiple things at once:
	// * selecting what command to run
	// * parsing the action from the command
	// * running the before/afterBake() method to be tested
	// * checking the results of the previous step
	//
	// The reason for this is to provide a DSL-like invocation for the
	// method, so that it is obvious at a glance what the method is testing;
	// an example call to this method may look like
	//
	//     assert(t, commandDirective, goodCmd, succeeds, afterBake)
	//
	// alternatives:
	// calls to smaller methods in calling method
	//     pros: dedicated methods
	//     cons: calls are duplicated in sequence in every method call
	//
	// calls to smaller methods in this method
	//     pros: dedicated methods
	//     cons: it is better to be able to control variables locally, so
	//           that if the command produces side-effects, they can check
	//           that the side effects occurred as expected within the same
	//           function
	//
	// chain method DSL (either "stateful" or "stateless"):
	//     pros: dedicated methods
	//           more DSL-like
	//     cons: more complicated

	// Arrange

	// cmdLine is enumerated so that the actual command to be run is local
	// to this assert method -- this is done so that changes to the command
	// to be run will affect all callers. It is also helpful to have the
	// command to be run is defined in the same method as it is run so that
	// expected side-effects can be asserted locally. In this way, cmdLine
	// is 'linked' to the outcomes at the bottom of this method.
	var cmdLine string
	switch cmd {
	case goodCmd:
		cmdLine = "true"
	case badCmd:
		cmdLine = "false"
	case errCmd:
		cmdLine = ""
	default:
		t.Fatalf("%d is not a valid cmd", cmd)
	}
	action := string(directive) + cmdLine
	command := parseAction(t, "descr\n"+action+"\n")

	// Act
	var result *result
	var err error
	if beforeBake {
		result, err = command.beforeBake()
	} else {
		result, err = command.afterBake()
	}

	// Assert

	// general rule - possible outcomes of any method:
	// err != nil && result != nil : ok -> err
	// err != nil && result == nil : ok -> err
	// err == nil && result != nil : ok -> result
	// err == nil && result == nil : fatal error

	when := "before"
	if !beforeBake {
		when = "after"
	}

	if expect == errs {
		// not &&'ing to previous condition so that it's clearer that
		// this test occurs when we're expecting an error
		if err == nil {
			if result == nil {
				t.Fatalf("err and result were nil when expecting error")
			}
			// inv: err == nil && result != nil

			t.Fatalf("expected an error %s bake, got '%s'", when,
				result.detail())
		}
	} else {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// inv: err == nil

		if result == nil {
			t.Fatalf("err and result were nil")
		}
		// inv: err == nil && result != nil

		successExpectation := expect == succeeds

		if result.success() != successExpectation {
			expected := "succeed"
			if !successExpectation {
				expected = "fail"
			}

			t.Fatalf("expected '%s' test to %s %s bake: %s", action,
				expected, when, result.detail())
		}
	}
}

func parseAction(t *testing.T, def string) testAction {
	buf := bytes.NewBufferString(def)
	in := strio.NewLineReader(buf)

	tests, err := readTypeTests(in)
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
