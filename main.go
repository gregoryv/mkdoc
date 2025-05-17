// Command main formats plain text files.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	txtfmt(os.Stderr, os.Stdout, os.Stdin)
}

func txtfmt(err, out io.Writer, in io.Reader) {
	w := &bytes.Buffer{}
	r := &bytes.Buffer{}
	io.Copy(r, in)
	next := func(step func()) {
		step()
		io.Copy(r, w)
	}

	// first pass; include files
	next(func() { incfile(w, r, "<>") })

	// parse links early
	var links map[string]string
	next(func() { links = parselinks(w, r) })

	next(func() { dropcomments(w, r) })

	next(func() { rfcindent(w, r) })

	// second pass; parse toc and index sections
	var toc bytes.Buffer
	cols := 69
	next(func() { parsetoc(w, &toc, r, cols) })

	// insert toc
	next(func() { inserttoc(w, r, &toc) })

	// replace links, also includes reference links
	next(func() { replacelinks(w, r, links) })

	fmt.Fprintln(out, `<!DOCTYPE html>

<meta charset="utf-8">
<pre>`)
	io.Copy(out, r)
}

type stepFn func()

var handleError = log.Fatal
