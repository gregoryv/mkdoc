package txtfmt

import (
	"bytes"
	"flag"
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

	// second pass; parse toc and index sections
	var toc bytes.Buffer
	next(func() { parsetoc(w, &toc, r, *cols) })

	// insert toc
	next(func() { inserttoc(w, r, &toc) })

	// parse links
	var links map[string]string
	next(func() { links = parselinks(w, r) })

	// replace links, also includes reference links
	next(func() { replacelinks(w, r, links) })

	io.Copy(c.out, r)
	return nil
}

type stepFn func()
