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
	condDelim = '?' // Denotes the start/end of a conditional insert
	condElsif = ':' // Denotes the else of a conditional insert
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

	err := copyUntilDelim(&in, out)
	for err == nil && !isEOF(&in) {
		if err = d.expandDirective(&in, out); err != nil {
			break
		}
		err = copyUntilDelim(&in, out)
	}
	if err != nil {
		return err
	}

	return out.Flush()
}

func isEOF(in *scanner.Scanner) bool {
	return in.Peek() == scanner.EOF
}

func copyUntilDelim(in *scanner.Scanner, out *bufio.Writer) error {
	var err error
	for err == nil && !isEOF(in) && in.Peek() != lDelim && in.Peek() != rDelim {
		err = copyNext(in, out)
	}
	return err
}

func copyNext(in *scanner.Scanner, out *bufio.Writer) error {
	c := in.Next()
	if out != nil {
		if n, err := out.WriteRune(c); err == nil && n < 1 {
			return fmt.Errorf("Couldn't write: %c", c)
		} else if err != nil {
			return err
		}
	}
	return nil
}

func (d *Dict) expandDirective(in *scanner.Scanner, out *bufio.Writer) error {
	var err error

	c := in.Next()
	switch c {
	case lDelim:
		switch in.Peek() {
		case lDelim:
			return copyNext(in, out)
		case condDelim:
			in.Next()
			err = d.expandCond(in, out)
		default:
			if err = d.expandVar(in, out); err == nil {
				err = match(in, rDelim)
			}
		}
	case rDelim:
		if err = writeString(out, "}"); err == nil {
			err = match(in, rDelim)
		}
	case scanner.EOF:
		err = parseErr(in, "Expected '%c' or '%c', got EOF", lDelim, rDelim)
	default:
		err = parseErr(in, "Expected '%c' or '%c', got '%c'", lDelim, rDelim, c)
	}

	return err
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

func parseErr(s *scanner.Scanner, msg string, params ...interface{}) error {
	p := s.Pos()
	text := fmt.Sprintf(msg, params...)
	return fmt.Errorf("%s[%d:%d] %s", p.Filename, p.Line, p.Column-1, text)
}

func (d *Dict) expandVar(in *scanner.Scanner, out *bufio.Writer) error {
	name, err := readVar(in)
	if err != nil {
		return err
	}

	val, ok := (*d)[name]
	if !ok {
		return parseErr(in, "Unknown variable '%s'", name)
	}

	return writeString(out, val)
}

func readVar(in *scanner.Scanner) (s string, err error) {
	var buf bytes.Buffer

	for isVarRune(in.Peek()) {
		buf.WriteRune(in.Next())
	}

	s = buf.String()
	if in.Peek() != rDelim {
		err = parseErr(in, "Unexpected character '%c'", in.Next())
	} else if len(s) == 0 {
		err = parseErr(in, "Empty variable")
	}
	return
}

// Is `r` a legal in a variable name?
func isVarRune(r rune) bool {
	return 'A' <= r && r <= 'Z' ||
		'a' <= r && r <= 'z' ||
		'0' <= r && r <= '9'
}

func writeString(out *bufio.Writer, s string) error {
	if out == nil {
		return nil
	}

	n := 0
	var err error
	if n, err = out.WriteString(s); err == nil && n != len(s) {
		err = fmt.Errorf("Only wrote %d characters of '%s'(%d)", n, s, len(s))
	}
	return err
}

func (d *Dict) expandCond(in *scanner.Scanner, writer *bufio.Writer) error {
	expanded := false

	expand, err := d.evalBool(in)
	if err != nil {
		return err
	}

	for err == nil && !isEOF(in) {
		var out *bufio.Writer
		if !expanded && expand {
			out = writer
			expanded = true
		}

		if err = copyUntilDelim(in, out); err != nil {
			break
		}

		if in.Peek() == lDelim {
			in.Next()

			switch in.Peek() {
			case lDelim:
				err = copyNext(in, out)
			case condDelim:
				in.Next()
				if in.Peek() == rDelim {
					in.Next()
					return nil
				} else {
					err = d.expandCond(in, out)
				}
			case condElsif:
				in.Next()
				if in.Peek() == rDelim {
					in.Next()
					expand = true
				} else {
					expand, err = d.evalBool(in)
				}
			default:
				err = d.expandVar(in, out)
				if out == nil || err == nil {
					err = match(in, rDelim)
				}
			}
		} else {
			err = d.expandDirective(in, out)
		}
	}

	return err
}

func (d *Dict) evalBool(in *scanner.Scanner) (bool, error) {
	if name, err := readVar(in); err != nil {
		return false, err
	} else {
		return d.hasVar(name), match(in, rDelim)
	}
}

func (d *Dict) hasVar(name string) bool {
	_, ok := (*d)[name]
	return ok
}
