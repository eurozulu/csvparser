package csvfile

import (
	"github.com/eurozulu/csvparser"
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
	file          *CSVFile
	headerIndex   int
	offset        int64
	length        int64
	columnIndexes []int
}

//goland:noinspection GoMaybeNil
func (r *RowSet) ColumnNames() []string {
	if r == nil || !r.HasColumnNames() {
		return nil
	}
	return r.columnNames().ColumnNames
}

func (r *RowSet) HasRows() bool {
	return r.length > 0
}

func (r *RowSet) Rows(rowCount int) ([]Row, error) {
	f, err := r.openFile()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lr := io.LimitReader(f, r.length)
	rows := make([]Row, 0, rowCount)
	p := csvparser.NewCsvParser(lr, r.file.Delimiter)
	filterCols := len(r.columnIndexes) > 0
	for p.Scan() {
		if rowCount >= 0 && len(rows) >= rowCount {
			return rows, p.Err()
		}
		row := p.Row()
		size := int64(r.rowLength(row))
		r.offset += size
		r.length -= size
		if filterCols {
			row = filterColumns(row, r.columnIndexes)
		}
		rows = append(rows, row)

	}
	// reached end of rows,
	r.offset = r.file.Length
	r.length = 0
	return rows, p.Err()
}

func (r *RowSet) HasColumnNames() bool {
	return r.headerIndex >= 0 && r.headerIndex < len(r.file.ColumnHeaders)
}

func (r *RowSet) rowLength(row Row) int {
	delimit := r.file.Delimiter
	return len(strings.Join(row, delimit.ColumnDelimiter)) + len(delimit.LineDelimiter)
}

func (r *RowSet) columnNames() *ColumnHeader {
	if !r.HasColumnNames() {
		return nil
	}
	return r.file.ColumnHeaders[r.headerIndex]
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

func filterColumns(row Row, colIndexes []int) Row {
	var result Row
	for _, colIndex := range colIndexes {
		var cell string
		if colIndex >= 0 && colIndex < len(row) {
			cell = row[colIndex]
		}
		result = append(result, cell)
	}
	return result
}

//goland:noinspection GoResourceLeak
func (r *RowSet) openFile() (io.ReadCloser, error) {
	if r.offset >= r.file.Length {
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
