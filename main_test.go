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
		_ = "Any argument SHOULD display usage"
		os.Args = []string{"", "--help"} // first arg is command name
		before := handleError
		defer func() { handleError = before }()

		handleError = func(v ...any) { t.Log(v) }
		main()
	})

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
		mkdoc(e, w, r)
	})

	t.Run("", func(t *testing.T) {
		_ = "By default HTML SHOULD be written to stdout"
		var e bytes.Buffer
		var w bytes.Buffer
		os.Chdir("docs")
		r := load("example.txt")

		mkdoc(&e, &w, r)
		golden.AssertWith(t, w.String(), "out.html")
		golden.AssertWith(t, e.String(), "err.html")
	})

	t.Run("", func(t *testing.T) {
		_ = "missing closing bracket in link SHOULD warn"
		r := strings.NewReader("... [text ")
		mkdoc(e, w, r)
	})

	t.Run("", func(t *testing.T) {
		req := "warn on missing file"
		var e bytes.Buffer
		r := strings.NewReader("<cat nosuch.txt>")

		mkdoc(&e, w, r)

		got := e.String()
		if got == "" {
			t.Error(fail(req, "empty"))
		}
	})

	t.Run("", func(t *testing.T) {
		req := "warn on missing section identifier"
		var e bytes.Buffer
		r := strings.NewReader("§ab")

		mkdoc(&e, w, r)

		got := e.String()
		if got == "" {
			t.Error(fail(req, "empty"))
		}
	})
	t.Run("", func(t *testing.T) {
		req := "already anchored links SHOULD be ignored"
		w := &bytes.Buffer{}
		txt := `... [<a href="#x">text</a>] .. `
		r := strings.NewReader(txt)
		mkdoc(e, w, r)
		if err := contains(w.String(), txt); err != nil {
			t.Error(fail(req, err))
		}
	})

	t.Run("", func(t *testing.T) {
		_ = "untagged requirements SHOULD warn"
		r := strings.NewReader(`
First SHOULD(#R1) be ok.

Second SHOULD fail.
This SHOULD
fail.
And this SHOULD
NOT succeed.
Duplicate SHOULD(#R1) also fail.
`)
		e := &bytes.Buffer{}
		mkdoc(e, w, r)

		got := e.String()
		err := contains(got, "line: 4", "line: 6", "line: 8", "line: 9")
		if err != nil {
			t.Error(err)
		}
	})
}

func Benchmark(b *testing.B) {
	d := ioutil.Discard
	os.Chdir("docs")
	r := load("example.txt")

	for b.Loop() {
		mkdoc(d, d, r)
	}
}

// ----------------------------------------

func load(filename string) io.Reader {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}
	return bytes.NewReader(data)
}

func fail(req string, err any) string {
	return fmt.Sprintln("\nreq:", req, "\nerr:", err)
}

func contains(txt string, substr ...string) error {
	for _, str := range substr {
		if !strings.Contains(txt, str) {
			return fmt.Errorf("missing %q\ngot: %s", str, txt)
		}
	}
	return nil
}

func compare(got, exp string) error {
	if got != exp {
		return fmt.Errorf("not equal\ngot:\n%s\nexp:\n%s", got, exp)
	}
	return nil
}
