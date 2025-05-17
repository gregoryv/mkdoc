package txtfmt

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func incfile(w io.Writer, r io.Reader, delim string) {
	s := bufio.NewScanner(r)
	prefix := "<incfile "
	for s.Scan() {
		line := s.Text()

		if strings.HasPrefix(line, prefix) {
			f := line[len(prefix):]
			f = strings.TrimRight(f, ">")
			f = strings.TrimSpace(f)
			fh, err := os.Open(f)
			if err != nil {
				log.Fatal(err)
			}
			io.Copy(w, fh)
			fh.Close()
		} else {
			fmt.Fprintln(w, line)
		}
	}
}
