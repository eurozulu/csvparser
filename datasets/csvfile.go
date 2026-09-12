package datasets

import (
	"github.com/eurozulu/csvparser"
	"github.com/eurozulu/csvparser/utils"
	"io"
	"os"
	"strings"
)

// CSVRowSampleSize defines the maximum number of rows used to detect a column header row(s).
// Search will scan a block of this size for a header row. If at least one found, it will continue
// to search the next block of rows
// until it finds a row block with no headers.
// i.e. This defines the maximum size of any datasets which preceed another dataset, if the following set is
// to be detected.
// Preceeding datasets are expected to be smaller Datasets, decribing the 'main' Dataset, which will usually be the
// last data set in a file.
var CSVRowSampleSize = 100

// CSVFile represents a single CSV text file containing one or more data sets.
// Each file may contain zero or more sets of column header names,
// each representing a distinct data set within the file.
// Column headers must be identifiers, Start with a letter and only contain letters numbers or "-_#"
// Each data set is represented as the column name strings, along with the file offset of where those
// headers appear in the file.
// Each data set offset is defined by the row following the header,
// until the next header row offset or the end of the file
// if no further header sets exist. (File Length)
// The File contains the delimiters found in the file.
// The column delimiter defaults to comma but will attempt to detect an alternative if a comma yields only one column.
// The data sets are read using the Column header offset.  This returns a Rowset of the rows from the given offset.
type CSVFile struct {
	Path      string                `csv:"path"`
	Delimiter *csvparser.Delimiters `csv:"delimiters"`
	Length    int64                 `csv:"length"`
	Modified  int64                 `csv:"modified"`
}

// DataSets returns all the Datasets found in the file.
func (f *CSVFile) DataSets() ([]*DataSet, error) {
	baseDS := &DataSet{
		CSVFile: *f,
		Offset:  0,
		Size:    f.Length,
		Columns: nil,
	}
	rowz, err := newRowIterator(baseDS)
	if err != nil {
		return nil, err
	}

	offset := int64(0)
	sets := []*DataSet{baseDS}
	for rowz.HasNextRows() {
		rows, err := rowz.NextRows(CSVRowSampleSize)
		if err != nil && err != io.EOF {
			return nil, err
		}

		found := false
		for _, row := range rows {
			if IsRowIdentifiers(row) {
				found = true
				if offset == 0 {
					// update baseDS with column name
					sets[0].Columns = ColumnNames(row)
				} else {
					set := &DataSet{
						CSVFile: *f,
						Offset:  offset,
						Size:    f.Length - offset,
						Columns: ColumnNames(row),
					}
					// update length of previous set to end at this offset
					last := len(sets) - 1
					sets[last].Size = offset - sets[last].Offset - 1
					sets = append(sets, set)
				}
			}
			offset += sizeOfRow(row, *f.Delimiter)
		}
		if !found {
			break
		}
	}

	return sets, nil
}

func IsRowIdentifiers(row Row) bool {
	if len(row) == 0 {
		return false
	}
	for _, s := range row {
		if s == "" {
			continue
		}
		if !utils.IsIdentifier(s, ColumnNameChars...) {
			return false
		}
	}
	return true
}

func sizeOfRow(row Row, delimit csvparser.Delimiters) int64 {
	offset := int64(len(strings.Join(row, delimit.ColumnDelimiter)))
	return offset + int64(len(delimit.LineDelimiter))
}

//goland:noinspection GoResourceLeak
func openFileBlock(path string, offset, length int64) (io.ReadCloser, error) {
	if length <= 0 {
		return nil, io.EOF
	}
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if offset > 0 {
		if _, err = r.Seek(offset, 1); err != nil {
			return nil, err
		}
	}
	return newLimitReaderCloser(r, length), nil
}

type limitReaderCloser struct {
	io.LimitedReader
}

func (l limitReaderCloser) Close() error {
	if rc, ok := l.R.(io.Closer); ok {
		return rc.Close()
	}
	return nil
}

func newLimitReaderCloser(r io.ReadCloser, limit int64) io.ReadCloser {
	return &limitReaderCloser{io.LimitedReader{
		R: r,
		N: limit,
	}}
}
