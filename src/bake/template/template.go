// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package template

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"text/scanner"
)

const (
	lDelim    = '{' // Denotes the start of a template directive
	rDelim    = '}' // Denotes the end of a template directive
)

type Dict map[string]string

func (d *Dict) ExpandStr(src string) (string, error) {
	var out bytes.Buffer
	in := bytes.NewBufferString(src)
	err := d.Expand(in, &out)

	if err != nil {
		return "", fmt.Errorf("'%s'%v", src, err)
	}

	return out.String(), nil
}

func (d *Dict) Expand(reader io.Reader, writer io.Writer) error {
	out := bufio.NewWriter(writer)
	var in scanner.Scanner
	in.Init(reader)

	err := d.expandText(&in, out)
	if err != nil {
		return err
	}

	if !isEOF(&in) {
		return parseErr(&in, "Unexpected closing statement: '%c'", in.Peek())
	}

	return out.Flush()
}

// Copy text until an escape sequence (EOF or directive) occurs.
//
// `out` may be `nil` to indicate that output is to be discarded, in which case
// errors are ignored.
func (d *Dict) expandText(in *scanner.Scanner, out *bufio.Writer) error {
	for finished := isEOF(in); !finished; finished = isEOF(in) {
		copyChar := false

		switch in.Peek() {
		case lDelim:
			in.Next()

			switch in.Peek() {
			case scanner.EOF:
				return parseErr(in,
					"Expected directive, got EOF")
			case lDelim:
				copyChar = true
			default:
				if err := d.expandVar(in, out); err != nil {
					return err
				}
			}
		case rDelim:
			in.Next()

			switch in.Peek() {
			case scanner.EOF:
				return parseErr(in, "Expected '%c', got EOF",
					rDelim)
			case rDelim:
				copyChar = true
			default:
				return parseErr(in, "Expected '%c', got '%c'",
					rDelim, in.Peek())
			}
		default:
			copyChar = true
		}

		if copyChar {
			c := in.Next()
			if out != nil {
				if n, err := out.WriteRune(c); err == nil && n < 1 {
					return fmt.Errorf("Couldn't write: %c", c)
				} else if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func isEOF(in *scanner.Scanner) bool {
	return in.Peek() == scanner.EOF
}

func parseErr(s *scanner.Scanner, msg string, params ...interface{}) error {
	p := s.Pos()
	text := fmt.Sprintf(msg, params...)
	return fmt.Errorf("%s[%d:%d] %s", p.Filename, p.Line, p.Column, text)
}

func (d *Dict) expandVar(in *scanner.Scanner, out *bufio.Writer) error {
	name, err := readVar(in)
	if err != nil {
		return err
	}

	if err = match(in, rDelim); err != nil {
		return err
	}

	if out != nil {
		val, ok := (*d)[name]
		if !ok {
			return parseErr(in, "Unknown variable '%s'", name)
		}

		if n, err := out.WriteString(val); n != len(val) {
			return parseErr(in, "Only wrote %d characters of '%s'(%d)",
				n, val, len(val))
		} else if err != nil {
			return err
		}
	}

	return out.Flush()
}

func readVar(in *scanner.Scanner) (s string, err error) {
	var buf bytes.Buffer

	for isVarRune(in.Peek()) {
		buf.WriteRune(in.Next())
	}

	s = buf.String()
	if len(s) == 0 {
		err = parseErr(in, "Unexpected character '%c'", in.Next())
	}

	return
}

// Is `r` a legal in a variable name?
func isVarRune(r rune) bool {
	return 'A' <= r && r <= 'Z' ||
		'a' <= r && r <= 'z' ||
		'0' <= r && r <= '9'
}

// Consume the next rune in `n` and return an error if it's not `r`.
func match(in *scanner.Scanner, r rune) error {
	var err error

	if c := in.Next(); c == scanner.EOF {
		err = parseErr(in, "Expected '%c', got EOF", r)
	} else if c != r {
		err = parseErr(in, "Expected '%c', got '%c'", r, c)
	}

	return err
}
