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

	// first pass include files
	var (
		first bytes.Buffer
		last  bytes.Buffer
	)

	io.Copy(&first, c.in)

	io.Copy(&last, &first)

	io.Copy(c.out, &last)
	return nil
}
