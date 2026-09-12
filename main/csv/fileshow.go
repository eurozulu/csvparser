package main

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/csvparser/csvfile"
	"strconv"
	"strings"
)

type fileShow struct {
}

func (fs fileShow) ShowFilesInfo(files []*csvfile.CSVFile) {
	for _, file := range files {
		fs.ShowFileInfo(file)
	}
}

func (fs fileShow) ShowFileInfo(file *csvfile.CSVFile) {
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
	if len(file.ColumnHeaders) == 0 {
		buf.WriteString("<no column names found>\n")
	} else {
		buf.WriteString("Column names:\n")
		for _, columnName := range file.ColumnHeaders {
			buf.WriteString("\t")
			buf.WriteString(strings.Join(columnName.ColumnNames, file.Delimiter.ColumnDelimiter))
			buf.WriteString("\n\t")
			buf.WriteString("offset: ")
			buf.WriteString(strconv.FormatInt(columnName.Offset, 10))
			buf.WriteString("\n")
		}
	}
	fmt.Println(buf.String())
}
