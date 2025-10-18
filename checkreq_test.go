package stp

import (
	"strings"
	"testing"
)

func TestCheckRequirements_duplicate(t *testing.T) {
	req := `requirement identifiers MUST not be reused`

	err, out, in, txt := input("Hepp MUST(#R1) here.\n ... \n Dong SHOULD(#R1) it.\n")

	CheckRequirements(err, out, in)

	got := format(err, out, txt)
	if !strings.Contains(got, "duplicate") {
		t.Error(got, req)
	}
}

func TestCheckRequirements_err(t *testing.T) {
	req := `requirement keywords MUST be followed by (#R\d)`

	cases := []string{
		"MUST NOT ...",
		"SHALL NOT ...",
		"SHOULD NOT ...",
		"MUST ...",
		"REQUIRED ...",
		"SHALL ...",
		"SHOULD ...",
		"RECOMMENDED ...",
		"... MAY ...",
		"OPTIONAL ...",
		// multiline
		"MUST\nNOT ...",
	}

	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			err, out, in, txt := input(c)

			CheckRequirements(err, out, in)

			got := format(err, out, txt)
			if err.Len() == 0 {
				t.Error(got, req)
			}
			if out.String() != c {
				t.Error(got, "stdout should match stdin")
			}
		})
	}
}

func TestCheckRequirements(t *testing.T) {
	req := `requirement keywords MUST be followed by (#R\d)`

	cases := []string{
		// ok
		`   "quoted MAY ..."`,
		"MUST NOT(#R1)",
		"SHALL NOT(#R2)",
		"SHOULD NOT(#R3)",
		"MUST(#R4)",
		"REQUIRED(#R5)",
		"SHALL(#R6)",
		"SHOULD(#R7)",
		"RECOMMENDED(#R8)",
		"... MAY(#R9) ...",
		"OPTIONAL(#R10)",
		// multiline
		"MUST\nNOT(#R11)",
		"MUST\n(#R12)",
	}

	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			err, out, in, txt := input(c)

			CheckRequirements(err, out, in)

			got := format(err, out, txt)
			if err.Len() > 0 {
				t.Error(got, req)
			}
			if out.String() != c {
				t.Error(got, "stdout should match stdin")
			}
		})
	}
}
