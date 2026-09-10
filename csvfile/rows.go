package csvfile

import "csvparser"

type Row []string

type Rows struct {
	ColumnNames *ColumnNames
	Size        int64
	rows        *csvparser.CsvParser
}

func (r Rows) Scan() bool {
	return r.rows.Scan()
}

func (r Rows) Row() Row {
	return r.rows.Row()
}

func (r Rows) FileLength() int64 {
	os := int64(0)
	if r.ColumnNames != nil {
		os = r.ColumnNames.Offset
	}
	return r.Size - os
}
