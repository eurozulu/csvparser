package csvfile

import (
	"slices"
	"strings"
	"unicode"

	"github.com/eurozulu/csvparser"
	"io"
	"os"
)

const sampleSize = 25

var ColumnNameChars = "_-#"

var CustomLineDelimiter = ""
var CustomColumnDelimiter = ""

func IsIdentifier(s string) bool {
	r := []rune(s)
	if len(r) == 0 {
		return false
	}
	if !unicode.IsLetter(r[0]) && !slices.Contains([]rune(ColumnNameChars), r[0]) {
		return false
	}
	idChars := []rune(ColumnNameChars)
	for _, c := range r[1:] {
		if !unicode.IsLetter(c) &&
			!unicode.IsNumber(c) &&
			!slices.Contains(idChars, c) {
			return false
		}
	}
	return true
}

func IsRowIdentifiers(row Row) bool {
	if len(row) == 0 {
		return false
	}
	for _, s := range row {
		if s == "" {
			continue
		}
		if !IsIdentifier(s) {
			return false
		}
	}
	return true
}

func findIdentifierRows(rows []Row) []int {
	var indexes []int
	for i, row := range rows {
		if !IsRowIdentifiers(row) {
			continue
		}
		indexes = append(indexes, i)
	}
	return indexes
}

func offsetOfRow(index int, rows []Row, delimit csvparser.Delimiters) int64 {
	offset := int64(0)
	for i, row := range rows {
		if i >= index {
			break
		}
		offset += int64(len(strings.Join(row, delimit.ColumnDelimiter)))
		offset += int64(len(delimit.LineDelimiter))
	}
	return offset
}

func sampleRows(path string, size int, delimit *csvparser.Delimiters) ([]Row, error) {
	if size < 1 {
		return nil, nil
	}
	delimit = csvparser.DelimiterOrDefault(delimit)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	p := csvparser.NewCsvParser(f, delimit)
	p.IgnoreEmptyRows = true
	if !p.Scan() {
		return nil, io.EOF
	}
	row := p.Row()
	if len(row) == 1 {
		de, ok := DetectColumnDelimter(strings.Join(row, delimit.ColumnDelimiter))
		if ok {
			delimit.ColumnDelimiter = de
			p.Delimiter.ColumnDelimiter = de
			row = p.Row()
		}
	}
	rows := []Row{row}
	for p.Scan() {
		if len(rows) >= size {
			break
		}
		rows = append(rows, p.Row())
	}
	return rows, nil
}

func ParseCSvFile(path string) (*CSVFile, error) {
	d := csvparser.DelimiterOrDefault()
	if CustomLineDelimiter != "" {
		d.LineDelimiter = CustomLineDelimiter
	}
	if CustomColumnDelimiter != "" {
		d.ColumnDelimiter = CustomColumnDelimiter
	}

	rows, err := sampleRows(path, sampleSize, d)
	if err != nil {
		return nil, err
	}

	file := &CSVFile{
		Path:      path,
		Delimiter: d,
	}
	if fi, err := os.Stat(file.Path); err == nil {
		file.Length = fi.Size()
	}

	idRows := findIdentifierRows(rows)
	for _, idRow := range idRows {
		file.ColumnHeaders = append(file.ColumnHeaders, &ColumnHeader{
			ColumnNames: ColumnNames(rows[idRow]),
			Offset:      offsetOfRow(idRow, rows, *d),
		})
	}
	return file, nil
}
