package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func Test(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_ = "No arguments SHOULD do nothing"
		os.Args = []string{""} // first arg is command name
		before := handleError
		defer func() { handleError = before }()

		handleError = func(v ...any) { t.Fatal(v) }
		main()
	})

	e := ioutil.Discard
	w := ioutil.Discard
	r := strings.NewReader("")

	t.Run("", func(t *testing.T) {
		_ = "Format empty input SHOULD do nothing"
		txtfmt(e, w, r)
	})

	t.Run("", func(t *testing.T) {
		_ = "By default HTML SHOULD be written to stdout"
		var e bytes.Buffer
		var w bytes.Buffer
		os.Chdir("testdata")
		r := load("example.txt")

		txtfmt(&e, &w, r)
		golden.AssertWith(t, w.String(), "out.html")
		golden.AssertWith(t, e.String(), "err.html")
	})

	t.Run("", func(t *testing.T) {
		_ = "missing closing bracket in link SHOULD warn"
		r := strings.NewReader("... [text ")
		txtfmt(e, w, r)
	})

	t.Run("", func(t *testing.T) {
		_ = "already anchored links SHOULD be ignored"
		r := strings.NewReader(`... [<a href="#x">text</a>] .. `)
		txtfmt(e, w, r)
	})
}

// ----------------------------------------

func load(filename string) io.Reader {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}
	return bytes.NewReader(data)
}
