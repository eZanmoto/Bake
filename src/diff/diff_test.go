// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package diff

import (
	"fmt"
	"testing"
)

var (
	changeTypeNames = map[ChangeType]string{
		Add:  "Add",
		Same: "Same",
		Rem:  "Rem",
	}
)

func TestChangeListIndexUnderflow(t *testing.T) {
	xs := []string{"a"}
	ys := []string{"b", "c"}

	cl := Diff(xs, ys)

	if _, _, err := cl.Get(-1); err == nil {
		t.Fatalf("Expected index underflow error")
	}
}

func TestChangeListIndexOverflow(t *testing.T) {
	xs := []string{"a"}
	ys := []string{"b", "c"}

	cl := Diff(xs, ys)

	if _, _, err := cl.Get(2); err != nil {
		t.Fatalf("Unexpected index error: %v", err)
	} else if _, _, err = cl.Get(3); err == nil {
		t.Fatalf("Expected index overflow error")
	}
}

func TestEmptyChanges(t *testing.T) {
	xs := []string{}
	ys := []string{}

	cl := Diff(xs, ys)

	assertNumChanges(t, 0, cl)
}

func assertNumChanges(t *testing.T, n int, cl ChangeList) {
	if cl.Len() != n {
		for i := 0; i < cl.Len(); i++ {
			ct_, line_, _ := cl.Get(i)
			fmt.Printf("[%s]%s\n", ct_.toString(), line_)
		}
		t.Fatalf("Expected %d changes, got %d", n, cl.Len())
	}
}

func TestAddAllLines(t *testing.T) {
	xs := []string{}
	ys := []string{"a", "b", "c"}

	cl := Diff(xs, ys)

	assertNumChanges(t, 3, cl)
	assertLineChange(t, cl, 0, Add, "a")
	assertLineChange(t, cl, 1, Add, "b")
	assertLineChange(t, cl, 2, Add, "c")
}

func assertLineChange(t *testing.T, cl ChangeList, n int, ct ChangeType,
	line string) {
	if ct_, line_, err := cl.Get(n); err != nil {
		t.Fatalf("Unexpected error %v", err)
	} else if ct_ != ct {
		t.Fatalf("Line %d should be '%s', got '%s'",
			n, ct.toString(), ct_.toString())
	} else if line != line_ {
		t.Fatalf("Expected '%s', got '%s'", line, line_)
	}
}

func (ct *ChangeType) toString() string {
	return changeTypeNames[*ct]
}

func TestRemoveAllLines(t *testing.T) {
	xs := []string{"a", "b", "c"}
	ys := []string{}

	cl := Diff(xs, ys)

	assertNumChanges(t, 3, cl)
	assertLineChange(t, cl, 0, Rem, "a")
}

func TestAllLinesSame(t *testing.T) {
	xs := []string{"a", "b", "c"}
	ys := []string{"a", "b", "c"}

	cl := Diff(xs, ys)

	assertNumChanges(t, 3, cl)
	assertLineChange(t, cl, 0, Same, "a")
	assertLineChange(t, cl, 1, Same, "b")
	assertLineChange(t, cl, 2, Same, "c")
}

func TestMoveBlankLineForward(t *testing.T) {
	xs := []string{"", "a", "b"}
	ys := []string{"a", "b", ""}

	cl := Diff(xs, ys)

	assertNumChanges(t, 4, cl)
	assertLineChange(t, cl, 0, Rem, "")
	assertLineChange(t, cl, 1, Same, "a")
	assertLineChange(t, cl, 2, Same, "b")
	assertLineChange(t, cl, 3, Add, "")
}

func TestMoveBlankLineBackward(t *testing.T) {
	xs := []string{"a", "b", ""}
	ys := []string{"", "a", "b"}

	cl := Diff(xs, ys)

	assertNumChanges(t, 4, cl)
	assertLineChange(t, cl, 0, Add, "")
	assertLineChange(t, cl, 1, Same, "a")
	assertLineChange(t, cl, 2, Same, "b")
	assertLineChange(t, cl, 3, Rem, "")
}

func TestAddLinesToCenter(t *testing.T) {
	xs := []string{"a", "d"}
	ys := []string{"a", "b", "c", "d"}

	cl := Diff(xs, ys)

	assertNumChanges(t, 4, cl)
	assertLineChange(t, cl, 0, Same, "a")
	assertLineChange(t, cl, 1, Add, "b")
	assertLineChange(t, cl, 2, Add, "c")
	assertLineChange(t, cl, 3, Same, "d")
}

func TestRemLinesFromCenter(t *testing.T) {
	xs := []string{"a", "b", "c", "d"}
	ys := []string{"a", "d"}

	cl := Diff(xs, ys)

	assertNumChanges(t, 4, cl)
	assertLineChange(t, cl, 0, Same, "a")
	assertLineChange(t, cl, 1, Rem, "b")
	assertLineChange(t, cl, 2, Rem, "c")
	assertLineChange(t, cl, 3, Same, "d")
}
