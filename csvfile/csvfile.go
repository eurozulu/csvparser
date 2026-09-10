package csvfile

import (
	"csvparser"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type CSVFile struct {
	Path        string               `csv:"path"`
	Delimiter   csvparser.Delimiters `csv:"delimiters"`
	ColumnNames []ColumnNames        `csv:"column_names"`
	Length      int64                `csv:"length"`
}

type ColumnNames struct {
	Offset      int64    `csv:"offset"`
	ColumnNames []string `csv:"columns"`
}

func (f CSVFile) Rows(offset int64) (*Rows, error) {
	file, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	if offset > 0 {
		if _, err = file.Seek(offset, 0); err != nil {
			return nil, err
		}
	}
	colIndex := f.indexColumnNamesForOffset(offset)
	r := &Rows{Size: f.Length}
	if colIndex != -1 {
		r.ColumnNames = &f.ColumnNames[colIndex]
		if colIndex+1 < len(f.ColumnNames) {
			r.Size = f.ColumnNames[colIndex+1].Offset - offset
		}
	}
	in := &CountReader{
		R:     file,
		Limit: r.Size,
	}
	r.rows = csvparser.NewCsvParser(in, &f.Delimiter)
	return r, nil
}

func (f CSVFile) indexColumnNamesForOffset(offset int64) int {
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
