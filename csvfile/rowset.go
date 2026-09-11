package csvfile

import (
	"csvparser"
	"fmt"
	"io"
	"os"
	"strings"
)

type Row []string

// RowSet represents a single data set within a CSV file.
// Rows returns the first data row at its offset, skipping any column header.
// Rows is an iterator of Rows, in that each call to Rows, moves the Rowset onto
// the next available rows in the set automatically.
// If the  call to rows requests a number exceeding the available rows, the full number available is returned
// and the Rowset will be empty
type RowSet struct {
	file         *CSVFile
	columnsIndex int
	offset       int64
	length       int64
}

func (r *RowSet) ColumnNames() []string {
	if !r.HasColumnNames() {
		return nil
	}
	return r.columnNames().ColumnNames
}

func (r *RowSet) HasRows() bool {
	return r.offset < r.length
}

func (r *RowSet) Rows(rowCount int) ([]Row, error) {
	zz, _ := os.ReadFile(r.file.Path)
	fmt.Println(string(zz))

	f, err := r.openFile()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lf := &ReadLimit{
		Limit: r.length,
		in:    f,
	}
	rows := make([]Row, 0, rowCount)
	p := csvparser.NewCsvParser(lf, r.file.Delimiter)
	for p.Scan() {
		if rowCount >= 0 && len(rows) >= rowCount {
			break
		}
		row := p.Row()
		size := int64(r.rowLength(row))
		rows = append(rows, row)
		r.offset += size
		r.length -= size
	}
	if p.Err() != nil {
		return nil, p.Err()
	}
	return rows, nil
}

func (r *RowSet) HasColumnNames() bool {
	return r.columnsIndex >= 0 && r.columnsIndex < len(r.file.ColumnNames)
}

func (r *RowSet) rowLength(row Row) int {
	delimit := r.file.Delimiter
	return len(strings.Join(row, delimit.ColumnDelimiter)) + len(delimit.LineDelimiter)
}

func (r *RowSet) columnNames() *ColumnHeader {
	if !r.HasColumnNames() {
		return nil
	}
	return r.file.ColumnNames[r.columnsIndex]
}

func (r *RowSet) skipColumnNames() error {
	names := r.columnNames()
	if names == nil || names.Offset != r.offset {
		return nil
	}
	if _, err := r.Rows(1); err != nil {
		return err
	}
	return nil
}

func (r RowSet) openFile() (io.ReadCloser, error) {
	if r.offset >= r.length {
		return nil, io.EOF
	}
	file, err := os.Open(r.file.Path)
	if err != nil {
		return nil, err
	}
	if r.offset > 0 {
		if _, err = file.Seek(r.offset, 0); err != nil {
			return nil, err
		}
	}
	return file, nil
}
