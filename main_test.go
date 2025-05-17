package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func Test(t *testing.T) {
	t.Run("", func(t *testing.T) {
		req := "No arguments SHOULD do nothing"
		os.Args = []string{""} // first arg is command name
		before := handleError
		defer func() { handleError = before }()

		handleError = func(v ...any) { t.Fatal(req, v) }
		main()
	})

	w := ioutil.Discard
	r := strings.NewReader("")

	t.Run("", func(t *testing.T) {
		req := "Format empty input SHOULD do nothing"

		err := txtfmt(w, r)
		if err != nil {
			t.Error(fail(req, err))
		}
	})

	t.Run("", func(t *testing.T) {
		req := "By default HTML SHOULD be written to stdout"
		var w bytes.Buffer
		os.Chdir("testdata")
		r := load("example.txt")

		err := txtfmt(&w, r)
		if err != nil {
			t.Fatal(fail(req, err))
		}
		golden.AssertWith(t, w.String(), "out.html")
	})

	t.Run("", func(t *testing.T) {
		req := "missing closing bracket in link SHOULD fail"
		r := strings.NewReader("... [text ")
		err := txtfmt(w, r)
		if err != nil {
			t.Error(fail(req, err))
		}
	})
}

// ----------------------------------------

func fail(req string, err error) string {
	return fmt.Sprintln("\nreq:", req, "\nerr:", err)
}

func load(filename string) io.Reader {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}
	return bytes.NewReader(data)
}
