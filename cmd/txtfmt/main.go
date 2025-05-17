// Command txtfmt formats plain text files.
package main

import (
	"io"
	"log"
	"os"

	"github.com/gregoryv/txtfmt"
)

func main() {
	cmd := txtfmt.NewCommand()
	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	err := cmd.Run(os.Args[1:]...)
	if err != nil {
		log.SetFlags(0)
		handleError(err)
	}
}

var (
	// define globals so we can override them in tests
	stdin       io.Reader = os.Stdin
	stdout      io.Writer = os.Stdout
	stderr      io.Writer = os.Stderr
	handleError           = log.Fatal
)
