// Command main formats plain text files.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

var usage = `Usage: mkdoc < input.txt

A text processing tool to generate RFC like software specifications
from plain text files.

Example:  https://gregoryv.github.io/mkdoc
`

func main() {
	log.SetFlags(0)
	if len(os.Args) > 1 {
		handleError(usage)
		return
	}

	mkdoc(os.Stderr, os.Stdout, os.Stdin)
}

func mkdoc(err, out io.Writer, in io.Reader) {
	log.SetOutput(err)
	w := &bytes.Buffer{}
	r := &bytes.Buffer{}
	io.Copy(r, in)
	next := func(step func()) {
		step()
		r, w = w, r
	}

	// first pass; include files
	next(func() { cat(w, r, "<>") })

	// check that all requirements are indexed with (#R...)
	next(func() { checkreq(err, w, r) })

	next(func() { replacerefs(w, r) })

	// parse links early
	var links map[string]string
	next(func() { links = parselinks(w, r) })

	next(func() { dropcomments(w, r) })

	next(func() { sentenceSpace(w, r) })
	next(func() { emptyLines(w, r) })

	next(func() { rfcindent(w, r) })

	// second pass; parse toc and index sections
	var toc bytes.Buffer
	cols := 69
	next(func() { parsetoc(w, &toc, r, cols) })
	next(func() { linksections(w, r) })

	// insert toc
	next(func() { inserttoc(w, r, &toc) })

	// before replacing ordinary links
	next(func() { replaceRequirements(w, r) })

	// replace links, also includes reference links
	next(func() { replacelinks(w, r, links) })

	next(func() { replaceSections(w, r) })

	fmt.Fprintln(out, htmlHeader)
	io.Copy(out, r)
}

type stepFn func()

var handleError = log.Fatal

const htmlHeader = `<!DOCTYPE html>

<meta charset="utf-8">
<pre>`
