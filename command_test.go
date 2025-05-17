package txtfmt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/gregoryv/golden"
)

func Test(t *testing.T) {
	t.Run("", func(t *testing.T) {
		req := "Format empty input SHOULD do nothing"
		cmd := NewCommand()

		err := cmd.Run()
		if err != nil {
			t.Error(fail(req, err))
		}
	})

	t.Run("", func(t *testing.T) {
		req := "Unknown option SHOULD fail"
		cmd := NewCommand()

		err := cmd.Run("-no-such")
		if err == nil {
			t.Error(req)
		}
	})

	t.Run("", func(t *testing.T) {
		req := "By default HTML SHOULD be written to stdout"
		cmd := NewCommand()
		var stdout bytes.Buffer
		cmd.SetOut(&stdout)
		cmd.SetIn(load("testdata/example.txt"))

		err := cmd.Run()
		if err != nil {
			t.Fatal(fail(req, err))
		}
		golden.AssertWith(t, stdout.String(), "testdata/out.html")
	})
}

// ----------------------------------------

func fail(req string, err error) string {
	return fmt.Sprintln("\nreq:", req, "\nerr:", err)
}

func compare(buf io.Writer, exp string) error {
	got := buf.(*bytes.Buffer).String()
	if got != exp {
		return fmt.Errorf("not equal\ngot:\n%s\nexp:\n%s", got, exp)
	}
	return nil
}

func load(filename string) io.Reader {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}
	return bytes.NewReader(data)
}
