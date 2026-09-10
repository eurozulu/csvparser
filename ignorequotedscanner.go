package csvparser

import (
	"bufio"
	"io"
	"strings"
	"bytes"
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

func doubleQuote(s string) string {
	if isDoubleQuoted(s) {
		return s
	}
	sb := &strings.Builder{}
	sb.WriteRune('"')
	sb.WriteString(escapeQuotes(s))
	sb.WriteRune('"')
	return sb.String()
}

func undoubleQuote(s string) string {
	if !isDoubleQuoted(s) {
		return s
	}
	return unescapeQuotes(strings.Trim(s, "\""))
}

func isDoubleQuoted(s string) bool {
	return strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")
}

func escapeQuotes(s string) string {
	sb := &strings.Builder{}
	for len(s) > 0 {
		i := strings.Index(s, "\"")
		if i == -1 {
			sb.WriteString(s)
			break
		}
		sb.WriteString(s[:i])
		if i > 0 && s[i-1] != '\\' {
			sb.WriteRune('\\')
		}
		sb.WriteString(s[i : i+1])
		s = s[i+1:]
	}
	return sb.String()
}

func unescapeQuotes(s string) string {
	return strings.ReplaceAll(s, "\\\"", "\"")
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
	return result
}
