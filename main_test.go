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

func Test_main(t *testing.T) {
	t.Run("", func(t *testing.T) {
		req := "No arguments SHOULD do nothing"
		os.Args = []string{""} // first arg is command name
		var failed bool
		handleError = func(v ...any) {
			t.Log(v)
			failed = true
		}

		main()
		if failed {
			t.Error(req)
		}
	})
}

func Test(t *testing.T) {
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
