package stp

import (
	"bytes"
	"os"
	"path"
	"strings"
	"testing"
)

func TestCat_nosuchfile(t *testing.T) {
	req := "<cat FILENAME> cat directive SHOULD fail if file doesn't exist"
	var err, in bytes.Buffer
	in.WriteString(`<cat nosuchfile>`)

	Cat(&err, nil, &in)

	if v := err.String(); !strings.Contains(v, "open nosuchfile") {
		t.Error(v, "\n", req)
	}
}

func TestCat(t *testing.T) {
	req := "<cat FILENAME> directive SHOULD include the named file"

	// generate temporary file
	filename := path.Join(t.TempDir(), "x")
	err := os.WriteFile(filename, []byte("world"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]string{
		"beginning": "<cat %s>\n...",
		"middle":    "...\n<cat %s>\n...",
		"end":       "...\n<cat %s>",
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err, out, in, txt := inputf(c, filename)

			Cat(err, out, in)

			if err.Len() > 0 {
				t.Error(err.String())
			}
			if v := out.String(); !strings.Contains(v, "world") {
				t.Error(format(err, out, txt), req)
			}
		})
	}
}
