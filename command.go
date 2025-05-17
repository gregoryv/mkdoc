package txtfmt

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// NewCommand returns a command ready for use.
func NewCommand() *Command {
	var cmd Command
	cmd.SetIn(strings.NewReader(""))
	cmd.SetOut(ioutil.Discard)
	cmd.SetErr(ioutil.Discard)
	return &cmd
}

type Command struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func (c *Command) Run(args ...string) error {
	fs := flag.NewFlagSet("txtfmt", flag.ContinueOnError)
	cols := fs.Int("cols", 69, "text width")

	fs.SetOutput(c.err)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	_ = cols

	w := &bytes.Buffer{}
	r := &bytes.Buffer{}
	io.Copy(r, c.in)
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
	next(func() { parsetoc(w, &toc, r, *cols) })

	// insert toc
	next(func() { inserttoc(w, r, &toc) })

	// replace links, also includes reference links
	next(func() { replacelinks(w, r, links) })

	fmt.Fprintln(c.out, `<!DOCTYPE html>

<meta charset="utf-8">
<pre>`)
	io.Copy(c.out, r)
	return nil
}

type stepFn func()
