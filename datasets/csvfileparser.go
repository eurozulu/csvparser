package datasets

import (
	"errors"
	"fmt"
	"github.com/eurozulu/csvparser"
	"os"
)

var ColumnNameChars = []rune("_-#")

func ParseCSvFile(path string, delimiter ...csvparser.Delimiters) (*CSVFile, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var delimit *csvparser.Delimiters
	if len(delimiter) > 0 {
		delimit = &delimiter[0]
	} else {
		delimit, err = DetectDelimiters(path)
		if err != nil {
			return nil, err
		}
	}
	return &CSVFile{
		Path:      path,
		Delimiter: delimit,
		Length:    fi.Size(),
		Modified:  fi.ModTime().Unix(),
	}, nil
}

func ParseCSVFiles(fileNames []string, delimiter ...csvparser.Delimiters) ([]*CSVFile, error) {
	var files []*CSVFile
	var errs []error
	for _, fileName := range fileNames {
		fz, err := ParseCSvFile(fileName, delimiter...)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to parse %q  %v", fileName, err))
			continue
		}
		files = append(files, fz)
	}
	return files, errors.Join(errs...)
}
