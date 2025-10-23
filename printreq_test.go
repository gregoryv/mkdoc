package stp

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_printLongRequirement(t *testing.T) {
	tc := testcase(`
When printing the list of requirements the text should
be aligned but to the right of the requirement id.

R?? text here
    and here

Max width of the text and requirement should be <= 69 characters
and words never split.
`)

	var buf bytes.Buffer
	id := "R1"
	line := "word word word word a."
	maxWidth := 15
	printLongRequirement(&buf, id, line, maxWidth)

	exp := ` word word
   word word a.`
	got := buf.String()
	if got != exp {
		t.Errorf("%s\n%q", tc, got)
	}
}

type testcase string

func (tc testcase) String() string {
	return fmt.Sprintf("%s%s%s", Yellow, string(tc), Reset)
}

const (
	// Reset
	Reset = "\x1b[0m"

	// Foreground (normal)
	Black   = "\x1b[30m"
	Red     = "\x1b[31m"
	Green   = "\x1b[32m"
	Yellow  = "\x1b[33m"
	Blue    = "\x1b[34m"
	Magenta = "\x1b[35m"
	Cyan    = "\x1b[36m"
	White   = "\x1b[37m"

	// Foreground (bright)
	BrightBlack   = "\x1b[90m"
	BrightRed     = "\x1b[91m"
	BrightGreen   = "\x1b[92m"
	BrightYellow  = "\x1b[93m"
	BrightBlue    = "\x1b[94m"
	BrightMagenta = "\x1b[95m"
	BrightCyan    = "\x1b[96m"
	BrightWhite   = "\x1b[97m"

	// Background (normal)
	BgBlack   = "\x1b[40m"
	BgRed     = "\x1b[41m"
	BgGreen   = "\x1b[42m"
	BgYellow  = "\x1b[43m"
	BgBlue    = "\x1b[44m"
	BgMagenta = "\x1b[45m"
	BgCyan    = "\x1b[46m"
	BgWhite   = "\x1b[47m"

	// Background (bright)
	BgBrightBlack   = "\x1b[100m"
	BgBrightRed     = "\x1b[101m"
	BgBrightGreen   = "\x1b[102m"
	BgBrightYellow  = "\x1b[103m"
	BgBrightBlue    = "\x1b[104m"
	BgBrightMagenta = "\x1b[105m"
	BgBrightCyan    = "\x1b[106m"
	BgBrightWhite   = "\x1b[107m"
)
