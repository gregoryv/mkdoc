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
	"text/template"
)

var tpl = template.Must(template.New("").Parse(
	`Usage: stp [OPTIONS]

A text processing tool to generate RFC like software specifications
from plain text files.

Example:  https://gregoryv.github.io/stp

{{.Options}}
Version..: {{.Version}}
Revision.: {{.ShortRevision}}
Author...: Gregory Vincic

`))

func usage() {
	var options bytes.Buffer
	fs.SetOutput(&options)
	fs.PrintDefaults()

	m := struct {
		Version       string
		ShortRevision string
		Options       string
	}{
		Version:       Version(),
		ShortRevision: Revision(6),
		Options:       options.String(),
	}
	tpl.Execute(stderr, m)
}

var (
	stderr io.Writer = os.Stderr
	stdout io.Writer = os.Stdout
	stdin  io.Reader = os.Stdin

	fs      = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	input   = fs.String("i", "", "File to read")
	output  = fs.String("o", "", "File to write")
	showVer = fs.Bool("v", false, "Show version and exit")
)

func main() {
	log.SetFlags(0)
	fs.Usage = usage

	// parse options
	err := fs.Parse(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		handleError(err)
		return
	}

	if *showVer {
		fmt.Fprintln(stdout, Version())
		return
	}

	var in io.Reader = stdin
	// use file input if given
	if *input != "" {
		fh, err := os.Open(*input)
		if err != nil {
			handleError(err)
		}
		in = fh
	}

	var out io.Writer = stdout
	// write to file if output is given
	if *output != "" {
		fh, err := os.Create(*output)
		if err != nil {
			handleError(err)
		}
		out = fh
	}

	run(stderr, out, in)
}

var handleError = log.Fatal
