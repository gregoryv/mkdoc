package main

import (
	"io/ioutil"
	"os"
	"testing"
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

	t.Run("", func(t *testing.T) {
		req := "Bad option SHOULD fail"
		os.Args = []string{"", "-no-such"}
		stderr = ioutil.Discard

		var failed bool
		handleError = func(v ...any) { failed = true }

		main()
		if !failed {
			t.Error(req)
		}
	})
}
