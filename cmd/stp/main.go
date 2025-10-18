// Command main formats plain text files.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/gregoryv/stp"
)

var tpl = template.Must(template.New("").Parse(
	`Usage: main [OPTIONS]

A text processing tool to generate RFC like software specifications
from plain text files.

Example:  https://gregoryv.github.io/main

{{.Options}}
Author...: Gregory Vincic

`))

func usage(w io.Writer) func() {
	return func() {
		var options bytes.Buffer
		fs.SetOutput(&options)
		fs.PrintDefaults()

		m := struct {
			Version       string
			ShortRevision string
			Options       string
		}{
			Version:       stp.Version(),
			ShortRevision: stp.Revision(6),
			Options:       options.String(),
		}
		tpl.Execute(w, m)
	}
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
	fs.Usage = usage(stderr)

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
		fmt.Fprintln(stdout, stp.Version())
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

var handleError = func(v ...any) {
	fmt.Fprint(stderr, v...)
	os.Exit(1)
}
