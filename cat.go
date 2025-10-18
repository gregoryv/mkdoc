package stp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Cat looks for <cat FILENAME> in injects the file.
func Cat(stderr, stdout io.Writer, stdin io.Reader) {
	r := bufio.NewReader(stdin)
	prefix := "<cat "

	for {
		line, err := r.ReadString('\n')
		if len(line) == 0 && errors.Is(err, io.EOF) {
			return
		}
		if strings.HasPrefix(line, prefix) {
			f := line[len(prefix):]
			f = strings.Trim(f, " >\n")
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprint(stderr, err)
				continue
			}
			io.Copy(stdout, fh)
			fh.Close()
		} else {
			fmt.Fprint(stdout, line)
			if errors.Is(err, io.EOF) {
				return
			}
		}
	}
}
