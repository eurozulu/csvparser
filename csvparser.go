package csvparser

import (
	"io"
)

type Parser interface {
	Scan() bool
	Row() []string
}

type CsvParser struct {
	Delimiter       *Delimiters
	IgnoreEmptyRows bool
	lines           *IgnoreQuotedScanner
}

func (r *CsvParser) Scan() bool {
	if !r.lines.Scan() {
		return false
	}
	if !r.IgnoreEmptyRows && len(r.lines.Bytes()) == 0 {
		return false
	}
	return true
}

func (r *CsvParser) Row() []string {
	if len(r.lines.Bytes()) == 0 {
		return nil
	}
	return IgnoreQuotedSplit(r.lines.Text(), r.Delimiter.ColumnDelimiter)
}

func (r *CsvParser) Err() error {
	return r.lines.Err()
}

func NewCsvParser(r io.Reader, delimiter ...*Delimiters) *CsvParser {
	delimit := DelimiterOrDefault(delimiter...)
	rwz := &CsvParser{
		Delimiter: delimit,
		lines:     NewIgnoreQuotedScanner(r, delimit.LineDelimiter),
	}
	return rwz
}
