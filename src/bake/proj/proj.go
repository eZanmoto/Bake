// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package proj provides project generation functionality.
package proj

import (
	"bake/template"
)

type Project struct {
	lang    string
	types   []string
	verbose bool
	dict    *template.Dict
}

func New(lg string, ts []string, v bool, vs map[string]string) Project {
	d := template.Dict{}
	for name, val := range vs {
		d[name] = val
	}
	for _, t := range ts {
		d[t] = ""
	}
	return Project{lg, ts, v, &d}
}

func (p *Project) IsOfType(t string) bool {
	for _, val := range p.types {
		if val == t {
			return true
		}
	}
	return false
}
