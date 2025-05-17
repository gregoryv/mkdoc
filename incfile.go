package txtfmt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	OPEN   byte = '<'
	CLOSE  byte = '>'
	delim       = string(OPEN) + string(CLOSE)
	indent      = ""
	tag         = "incfile"
	prefix      = []byte(tag + " ")
)

func incfile(w io.Writer, txt io.Reader, delim string) {
	if len(delim) == 2 {
		OPEN = delim[0]
		CLOSE = delim[1]
	}

	r := bufio.NewReader(txt)
	var err error
	next := openTag
	for {
		next, err = next(w, r)
		if err != nil {
			break
		}
	}
}

func openTag(w io.Writer, r *bufio.Reader) (parseFn, error) {
	// look for open character
	head, err := r.ReadBytes(OPEN)
	if err != nil {
		w.Write(head)
		return nil, err
	}

	// write everything before OPEN
	if l := len(head); l > 1 {
		w.Write(head[:l-1])
		head = head[l-1:]
	}

	tail, err := r.ReadBytes(CLOSE)

	// ignore other than our prefix
	if !bytes.HasPrefix(bytes.TrimSpace(tail), prefix) {
		w.Write(head)
		w.Write(tail)
		return openTag, nil // start over
	}

	// parse filename and arguments
	tail = bytes.TrimSpace(tail)
	filename := tail[len(prefix) : len(tail)-1]

	// include the file
	fh, err := os.Open(string(filename))
	if err != nil {
		handleError("incfile: ", err)
		return openTag, nil
	}
	s := bufio.NewScanner(fh)
	for s.Scan() {
		fmt.Fprint(w, indent, s.Text(), "\n")
	}
	fh.Close()

	return openTag, nil
}

type parseFn func(io.Writer, *bufio.Reader) (parseFn, error)

var handleError = log.Fatal
