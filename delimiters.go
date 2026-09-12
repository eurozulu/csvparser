package csvparser

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/csvparser/utils"
	"strings"
)

// Delimiters represent the delimiters used in the CSV file.
type Delimiters struct {
	LineDelimiter   string
	ColumnDelimiter string
}

var defaultDelimiters = Delimiters{
	LineDelimiter:   "\n",
	ColumnDelimiter: ",",
}

func DelimiterOrDefault(delimiters ...Delimiters) Delimiters {
	if len(delimiters) == 0 {
		return defaultDelimiters
	}
	return delimiters[0]
}

func (d *Delimiters) String() string {
	ld := d.LineDelimiter
	if strings.Contains(ld, " ") {
		ld = utils.DoubleQuote(ld)
	}
	buf := bytes.NewBufferString(ld)
	buf.WriteRune(' ')
	cd := d.ColumnDelimiter
	if strings.Contains(cd, " ") {
		cd = utils.DoubleQuote(cd)
	}
	buf.WriteString(cd)
	return buf.String()
}

func (d *Delimiters) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

func (d *Delimiters) UnmarshalText(text []byte) error {
	iqs := NewIgnoreQuotedScanner(bytes.NewReader(text), " ")
	var segs []string
	for iqs.Scan() {
		v := utils.UnescapeQuotes(iqs.Text())
		segs = append(segs, v)
		if len(segs) > 2 {
			return fmt.Errorf("invalid delimiter string %q", v)
		}
	}
	if len(segs) != 2 {
		return fmt.Errorf("invalid delimiter string %s", string(text))
	}

	d.LineDelimiter = strings.Trim(segs[0], "\" \t")
	d.ColumnDelimiter = strings.Trim(segs[1], "\" \t")
	return nil
}
