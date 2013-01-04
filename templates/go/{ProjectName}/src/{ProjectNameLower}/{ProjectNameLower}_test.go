// Copyright {Year} {Owner}. All rights reserved.
{?License:
// Use of this source code is governed by a {License}
// license that can be found in the LICENSE file.
}

package main

import (
	"testing"
)

func BasicTest(t *testing.T) {{
	if 1 == 2 {{
		t.Fatalf("Something is very wrong")
	}}
}}
