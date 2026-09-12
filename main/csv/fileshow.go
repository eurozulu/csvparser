package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/eurozulu/csvparser/datasets"
	"strconv"
	"strings"
)

type fileShow struct {
}

func (fs fileShow) ShowFilesInfo(files []*datasets.CSVFile) error {
	var errs []error
	for _, file := range files {
		if err := fs.ShowFileInfo(file); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (fs fileShow) ShowFileInfo(file *datasets.CSVFile) error {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("Path: ")
	buf.WriteString(file.Path)
	buf.WriteString("\n")
	buf.WriteString("Line/Column delimiters: ")
	buf.WriteString(strconv.Quote(file.Delimiter.String()))
	buf.WriteString("\n")
	buf.WriteString("length: ")
	buf.WriteString(strconv.FormatInt(file.Length, 10))
	buf.WriteString("\n")
	sets, err := file.DataSets()
	if err != nil {
		return err
	}
	if len(sets) == 0 {
		buf.WriteString("<no column names found>\n")
	} else {
		buf.WriteString("Column names:\n")
		for _, set := range sets {
			buf.WriteString("\t")
			buf.WriteString(strings.Join(set.Columns, file.Delimiter.ColumnDelimiter))
			buf.WriteString("\n\t")
			buf.WriteString("offset: ")
			buf.WriteString(strconv.FormatInt(set.Offset, 10))
			buf.WriteString("\n")
			buf.WriteString("size: ")
			buf.WriteString(strconv.FormatInt(set.Size, 10))
			buf.WriteString("\n")
		}
	}
	fmt.Println(buf.String())
	return nil
}
