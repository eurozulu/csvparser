package csvfile

import (
	"errors"
	"fmt"
	"github.com/eurozulu/csvparser"
	"path/filepath"
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
	Path        string                `csv:"path"`
	Delimiter   *csvparser.Delimiters `csv:"delimiters"`
	ColumnNames []*ColumnHeader       `csv:"column_names"`
	Length      int64                 `csv:"length"`
}

type ColumnHeader struct {
	Offset      int64    `csv:"offset"`
	ColumnNames []string `csv:"columns"`
}

func (f *CSVFile) RowSet(offset int64) (*RowSet, error) {
	colIndex := f.indexColumnNamesForOffset(offset)
	rs := &RowSet{
		file:         f,
		offset:       offset,
		length:       f.Length - offset, // initially grab to EOF, unless another column set follows.
		columnsIndex: colIndex,
	}

	if colIndex != -1 && colIndex+1 < len(f.ColumnNames) {
		rs.length = f.ColumnNames[colIndex+1].Offset - offset
	}
	return rs, rs.skipColumnNames()
}

func (f *CSVFile) indexColumnNamesForOffset(offset int64) int {
	for i := len(f.ColumnNames) - 1; i >= 0; i-- {
		if f.ColumnNames[i].Offset > offset {
			continue
		}
		return i
	}
	return -1
}

func ParseCSVFiles(pattern string) ([]*CSVFile, error) {
	var files []*CSVFile
	var errs []error
	fileNames, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
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
