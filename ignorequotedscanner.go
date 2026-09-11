package csvparser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type IgnoreQuotedScanner struct {
	delimit     string
	delimitSize int
	scn         *bufio.Scanner
}

func (sc *IgnoreQuotedScanner) Scan() bool {
	return sc.scn.Scan()
}

func (sc *IgnoreQuotedScanner) Bytes() []byte {
	return sc.scn.Bytes()
}

func (sc *IgnoreQuotedScanner) Text() string {
	return sc.scn.Text()
}

func (sc *IgnoreQuotedScanner) Err() error {
	return sc.scn.Err()
}

func (sc *IgnoreQuotedScanner) splitLinesFunc(data []byte, eof bool) (advance int, token []byte, err error) {
	if eof && len(data) == 0 {
		return 0, nil, nil
	}

	inQuotes := false
	for i := 0; i < len(data); i++ {
		switch data[i] {
		case '"':
			inQuotes = !inQuotes
		case sc.delimit[0]:
			l := i + sc.delimitSize
			if inQuotes || l > len(data) || string(data[i:l]) != sc.delimit {
				continue
			}
			return l, data[:i], nil
		}
	}
	if eof {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func NewIgnoreQuotedScanner(r io.Reader, delimiter string) *IgnoreQuotedScanner {
	rwz := &IgnoreQuotedScanner{
		delimit:     delimiter,
		delimitSize: len(delimiter),
		scn:         bufio.NewScanner(r),
	}
	rwz.scn.Split(rwz.splitLinesFunc)
	return rwz
}

func IgnoreQuotedSplit(data string, delimiter string) []string {
	return IgnoreQuotedSplitN(data, delimiter, -1)
}

func IgnoreQuotedSplitN(data string, delimiter string, n int) []string {
	scn := NewIgnoreQuotedScanner(bytes.NewReader([]byte(data)), delimiter)
	var result []string
	unlimit := n < 0
	for scn.Scan() {
		if !unlimit && len(result) >= n {
			break
		}
		result = append(result, scn.Text())
	}
	if strings.HasSuffix(data, delimiter) {
		result = append(result, scn.Text())
	}
	return result
}
