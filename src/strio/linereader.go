// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package strio

type LineReader interface {
	// Reads the next line from the underlying stream.
	//
	// If the returned error is io.EOF, process the line as normal, e.g. if
	// stream contains "a\nb", calls to ReadLine() will return the
	// following:
	//
	//     ReadLine() -> ("a", nil)
	//     ReadLine() -> ("b", io.EOF)
	//
	// If the stream contains "a\nb\n", the calls will return the following:
	//
	//     ReadLine() -> ("a", nil)
	//     ReadLine() -> ("b", nil)
	//     ReadLine() -> ("", io.EOF)
	ReadLine() (string, error)

	// Same as ReadLine() but removes the trailing newline if present.
	ChompLine() (string, error)
}
