package stp

import (
	"bytes"
	"fmt"
)

func inputf(format string, v ...any) (err, out, in *bytes.Buffer, txt string) {
	return input(fmt.Sprintf(format, v...))
}

func input(v string) (err, out, in *bytes.Buffer, txt string) {
	txt = v
	err = &bytes.Buffer{}
	out = &bytes.Buffer{}
	in = bytes.NewBufferString(v)
	return
}

func format(err, out *bytes.Buffer, txt string) string {
	return fmt.Sprintf(`
----- stdin  -----
%s
----- stdout -----
%s
----- stderr -----
%s
------------------
`, txt, out.String(), err.String())
}
