package datasets

import (
	"crypto/sha1"
	"io"
)

// DataSet represents a subset of rows from a csvfile.
// Each CSV file may contain multiple tables, the data sets breaks down the rows in the CSV
// into discreet datasets, with their own column names and position data.
type DataSet struct {
	CSVFile
	Offset  int64       `csv:"offset"`
	Size    int64       `csv:"size"`
	Columns ColumnNames `csv:"column-names"`
}

func (ds *DataSet) Id() ([]byte, error) {
	r, err := openFileBlock(ds.Path, ds.Offset, ds.Size)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (ds *DataSet) Rows(columnNames ...string) (*RowIterator, error) {
	return newRowIterator(ds, columnNames...)
}
