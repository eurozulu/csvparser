package csvfile

import (
	"errors"
	"fmt"
	"github.com/eurozulu/csvparser"
	"github.com/eurozulu/csvparser/utils"
	"path/filepath"
	"slices"
	"strings"
)

// CSVFile represents a single CSV text file containing one or more data sets.
// Each file may contain zero or more sets of column header names,
// each representing a distinct data set within the file.
// Column headers must be identifiers, Start with a letter and only contain letters numbers or "-_#"
// Each header set is represented as the column name strings,
// along with the file offset of where those headers appear in the file.
// Each data set is defined by the row following the header, until the next header row offset or the end of the file
// if no further header sets exist. (File Length)
// The File contains the delimiters found in the file.  The line delimiter is assumed to be new line.
// The column delimiter defaults to comma but will attempt to detect an alternative if a comma yields only one column.
// The data sets are read using the Column header offset.  This returns a Rowset of the rows from the given offset.
type CSVFile struct {
	Path          string                `csv:"path"`
	Delimiter     *csvparser.Delimiters `csv:"delimiters"`
	ColumnHeaders []*ColumnHeader       `csv:"column-headers"`
	Length        int64                 `csv:"length"`
	Modified      int64                 `csv:"modified"`
}

type ColumnHeader struct {
	Offset      int64       `csv:"offset"`
	ColumnNames ColumnNames `csv:"columns"`
}

func (f *CSVFile) RowSet(columnNames ...string) (*RowSet, error) {
	headIndex := f.indexOfHeaderByNames(columnNames)
	if headIndex == -1 && len(columnNames) > 0 {
		return nil, fmt.Errorf("column names %q not known", strings.Join(columnNames, ", "))
	}

	offset := int64(0)
	var colIndexes []int
	if headIndex > -1 {
		offset = f.ColumnHeaders[headIndex].Offset
		if len(columnNames) != len(f.ColumnHeaders[headIndex].ColumnNames) {
			colIndexes = f.indexesOfColumnNames(f.ColumnHeaders[headIndex].ColumnNames, columnNames)
		}
	}

	rs := &RowSet{
		file:          f,
		offset:        offset,
		length:        f.Length - offset, // initially grab to EOF, unless another column set follows.
		headerIndex:   headIndex,
		columnIndexes: colIndexes,
	}

	if headIndex != -1 && headIndex+1 < len(f.ColumnHeaders) {
		rs.length = f.ColumnHeaders[headIndex+1].Offset - offset
	}
	return rs, rs.skipColumnNames()
}

func (f *CSVFile) indexOfHeaderByNames(names ColumnNames) int {
	for i := len(f.ColumnHeaders) - 1; i >= 0; i-- {
		if !utils.ContainsAll(f.ColumnHeaders[i].ColumnNames, names) {
			continue
		}
		return i
	}
	return -1
}

func (f *CSVFile) indexesOfColumnNames(headNames, colNames []string) []int {
	indexes := make([]int, len(colNames))
	for i, name := range colNames {
		indexes[i] = slices.Index(headNames, name)
	}
	return indexes
}

func ParseCSVFiles(pattern string) ([]*CSVFile, error) {
	var files []*CSVFile
	var errs []error
	fileNames, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(fileNames) == 0 {
		return nil, fmt.Errorf("no files found in %s", pattern)
	}
	for _, fileName := range fileNames {
		fz, err := ParseCSvFile(fileName)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to parse %q  %v", fileName, err))
			continue
		}
		files = append(files, fz)
	}
	return files, errors.Join(errs...)
}
