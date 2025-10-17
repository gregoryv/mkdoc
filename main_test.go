package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func Test(t *testing.T) {
	t.Run("", func(t *testing.T) {
		req := "Any argument SHOULD display usage"
		os.Args = []string{"", "--unknown"} // first arg is command name

		catch(t)
		var buf bytes.Buffer

		// set globals, these will be reset by catch
		stderr = &buf
		handleError = func(_ ...any) {
			v := buf.String()
			if !strings.Contains(v, "Usage:") {
				t.Error(v, "\nFAIL:", req)
			}
		}
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

	e := io.Discard
	w := io.Discard
	r := strings.NewReader("")

	t.Run("", func(t *testing.T) {
		_ = "Format empty input SHOULD do nothing"
		run(e, w, r)
	})

	t.Run("", func(t *testing.T) {
		_ = "By default HTML SHOULD be written to stdout"
		var e bytes.Buffer
		var w bytes.Buffer
		os.Chdir("docs")
		r := load("example.txt")

		run(&e, &w, r)
		golden.AssertWith(t, w.String(), "out.html")
		golden.AssertWith(t, e.String(), "err.html")
	})

	t.Run("", func(t *testing.T) {
		_ = "missing closing bracket in link SHOULD warn"
		r := strings.NewReader("... [text ")
		run(e, w, r)
	})

	t.Run("", func(t *testing.T) {
		req := "warn on missing file"
		var e bytes.Buffer
		r := strings.NewReader("<cat nosuch.txt>")

		run(&e, w, r)

		got := e.String()
		if got == "" {
			t.Error(fail(req, "empty"))
		}
	})

	t.Run("", func(t *testing.T) {
		req := "warn on missing section identifier"
		var e bytes.Buffer
		r := strings.NewReader("§ab")

		run(&e, w, r)

		got := e.String()
		if got == "" {
			t.Error(fail(req, "empty"))
		}
	})

	t.Run("", func(t *testing.T) {
		req := "warn on missing double spaces after end of sentence"
		var e bytes.Buffer
		r := strings.NewReader("Hello you. What is up?")

		run(&e, w, r)

		got := e.String()
		if err := contains(got, "missing double space"); err != nil {
			t.Error(fail(req, got))
		}
	})

	t.Run("", func(t *testing.T) {
		req := "abbreviation followed by uppercase word is not a new sentence"
		cases := []string{
			"i.e. Uppercase word?",
			"i.e., Uppercase word?",
			"eg. Uppercase word?",
			"-d '{ ... OPTIONS ...}' http://www.example.com",
		}
		for _, v := range cases {
			var e bytes.Buffer
			r := strings.NewReader(v)

			run(&e, w, r)

			got := e.String()
			if err := contains(got, "missing double space"); err == nil {
				t.Error(fail(req, got))
			}
		}
	})
	t.Run("", func(t *testing.T) {
		req := "already anchored links SHOULD be ignored"
		w := &bytes.Buffer{}
		txt := `... [<a href="#x">text</a>] .. `
		r := strings.NewReader(txt)
		run(e, w, r)
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
		run(e, w, r)

		got := e.String()
		err := contains(got, "line: 4", "line: 6", "line: 8", "line: 9")
		if err != nil {
			t.Error(err)
		}
	})
}

func Benchmark(b *testing.B) {
	d := io.Discard
	os.Chdir("docs")
	r := load("example.txt")

	for b.Loop() {
		run(d, d, r)
	}
}

// ----------------------------------------

func catch(t *testing.T) {
	bstderr := stderr
	bhandleError := handleError
	t.Cleanup(func() {
		handleError = bhandleError
		stderr = bstderr
	})
}

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
