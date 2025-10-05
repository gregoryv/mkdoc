// Command main formats plain text files.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func usage() {
	fmt.Fprint(fs.Output(), `Usage: mkdoc [OPTIONS]

A text processing tool to generate RFC like software specifications
from plain text files.

Example:  https://gregoryv.github.io/mkdoc
`)
	fs.PrintDefaults()
}

var (
	fs    = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	input = fs.String("i", "", "Default is stdin")
)

func main() {
	log.SetFlags(0)
	fs.Usage = usage

	err := fs.Parse(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		handleError(err)
		return
	}

	if *input == "" && len(os.Args) > 1 {
		handleError(usage)
		return
	}

	var in io.Reader = os.Stdin
	if *input != "" {
		fh, err := os.Open(*input)
		if err != nil {
			handleError(err)
		}
		in = fh
	}

	mkdoc(os.Stderr, os.Stdout, in)
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

	// parse links early
	var links map[string]string
	next(func() { links = parselinks(w, r) })

	var requirements []string
	next(func() { requirements = parsereq(w, r) })

	// requirements must be indexed (#R...)
	next(func() { checkreq(err, w, r) }) // #R8
	next(func() { sentenceSpace(w, r) })
	next(func() { emptyLines(w, r) })
	next(func() { includeReq(w, r, requirements) })

	// lines starting with `[\d+] ...`
	next(func() { replacerefs(w, r) })

	next(func() { dropcomments(w, r) })

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
