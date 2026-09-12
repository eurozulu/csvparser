package datasets

import (
	"fmt"
	"github.com/eurozulu/csvparser"
	"io"
	"slices"
)

type Row []string

// RowIterator represents a stream of rows in a Dataset
type RowIterator struct {
	dataSet       *DataSet
	offset        int64
	length        int64
	columnIndexes []int
}

func (r *RowIterator) HasColumnNames() bool {
	return len(r.dataSet.Columns) > 0
}

func (r *RowIterator) Reset() error {
	r.offset = r.dataSet.Offset
	r.length = r.dataSet.Size
	return r.skipColumnNames()
}

func (r *RowIterator) ColumnNames() []string {
	if r == nil || !r.HasColumnNames() {
		return nil
	}
	names := r.dataSet.Columns
	if len(r.columnIndexes) == 0 {
		return names
	}
	return filterColumns(names, r.columnIndexes)
}

func (r *RowIterator) HasNextRows() bool {
	return r.length > int64(len(r.dataSet.Delimiter.LineDelimiter))
}

func (r *RowIterator) NextRows(rowCount int) ([]Row, error) {
	f, err := openFileBlock(r.dataSet.Path, r.offset, r.length)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	delimiters := *r.dataSet.Delimiter
	rows := make([]Row, 0, rowCount)
	p := csvparser.NewCsvParser(f, delimiters)
	p.IgnoreEmptyRows = true
	filterCols := len(r.columnIndexes) > 0
	for p.Scan() {
		if rowCount >= 0 && len(rows) >= rowCount {
			return rows, p.Err()
		}
		row := p.Row()
		size := sizeOfRow(row, delimiters)
		r.offset += size
		r.length -= size
		if filterCols {
			row = filterColumns(row, r.columnIndexes)
		}
		rows = append(rows, row)
	}
	// reached end of rows, before rowCount
	//r.offset = r.dataSet.Size
	r.length = 0
	return rows, p.Err()
}

func (r *RowIterator) SkipRows(rowCount int) error {
	if rowCount == 0 {
		return nil
	}
	if rowCount < 0 {
		return fmt.Errorf("can not skip negative amount of rows")
	}

	f, err := openFileBlock(r.dataSet.Path, r.offset, r.length)
	if err != nil {
		return err
	}
	defer f.Close()

	delimiters := r.dataSet.Delimiter
	ldLen := len(delimiters.LineDelimiter)
	scn := csvparser.NewIgnoreQuotedScanner(f, delimiters.LineDelimiter)
	count := 0
	for scn.Scan() {
		if count >= rowCount {
			return scn.Err()
		}
		size := int64(len(scn.Text()) + ldLen)
		r.offset += size
		r.length -= size
		count++
	}
	return io.EOF
}

func (r *RowIterator) skipColumnNames() error {
	if !r.HasColumnNames() || r.offset != r.dataSet.Offset {
		return nil
	}
	return r.SkipRows(1)
}

func columnNameIndexes(haveNames, wantNames []string) ([]int, error) {
	var found []int
	for _, name := range wantNames {
		i := slices.Index(haveNames, name)
		if i < 0 {
			return nil, fmt.Errorf("column name %q not found in %q", name, haveNames)
		}
		found = append(found, i)
	}
	return found, nil
}

func filterColumns[S ~string](row []S, colIndexes []int) []S {
	result := make([]S, len(colIndexes))
	for i, colIndex := range colIndexes {
		var cell S
		if colIndex >= 0 && colIndex < len(row) {
			cell = row[colIndex]
		}
		result[i] = cell
	}
	return result
}

func newRowIterator(dataSet *DataSet, columnNames ...string) (*RowIterator, error) {
	var indexes []int
	if len(columnNames) > 0 && columnNames[0] != "" {
		var err error
		indexes, err = columnNameIndexes(dataSet.Columns, columnNames)
		if err != nil {
			return nil, err
		}
	}

	ri := &RowIterator{
		dataSet:       dataSet,
		columnIndexes: indexes,
	}
	if err := ri.Reset(); err != nil {
		return nil, err
	}
	return ri, nil
}
