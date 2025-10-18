package stp

import "testing"

func TestCheckRequirements(t *testing.T) {
	cases := []string{
		// ok
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
	}

	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			req := "requirement keywords MUST be followed by (#R\\d)"
			err, out, in, txt := input(c)

			Cat(err, out, in)

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
