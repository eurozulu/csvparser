package main

import (
	"os"
	"fmt"
	"csvparser/csvfile"
	"bytes"
	"strconv"
	"strings"
	"path/filepath"
	"errors"
)

func main() {
	args := Arguments(os.Args[1:])
	fileNames := args.Parameters()
	if len(fileNames) == 0 {
		exitError(fmt.Errorf("must provide a filename of the csv file to parse"))
	}
	if flg, ok := args.Flag("delimitline"); ok {
		if flg == "" {
			exitError(fmt.Errorf("must provide a line delimiter"))
		}
		csvfile.CustomLineDelimiter = flg
	}
	if flg, ok := args.Flag("delimitcol"); ok {
		if flg == "" {
			exitError(fmt.Errorf("must provide a column delimiter"))
		}
		csvfile.CustomColumnDelimiter = flg
	}

	filez, err := parseFiles(fileNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		if len(filez) == 0 {
			return
		}
	}

	if _, ok := args.Flag("info"); ok {
		showFilesInfo(filez)
		return
	}
	_, showColHead := args.Flag("showheaders")

	if err = showFilesData(filez, showColHead); err != nil {
		exitError(err)
	}
}

func showFilesData(files []*csvfile.CSVFile, showHead bool) error {
	for _, f := range files {
		fmt.Println(f.Path)
		if err := showFileRows(f, showHead); err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}

func showFileRows(f *csvfile.CSVFile, showHead bool) error {
	offset := int64(0)
	for {
		rows, err := f.Rows(offset)
		if err != nil {
			return err
		}
		if !showHead && rows.ColumnNames != nil && offset == rows.ColumnNames.Offset {
			rows.Scan()
		}
		for rows.Scan() {
			fmt.Println(strings.Join(rows.Row(), f.Delimiter.ColumnDelimiter))
		}
		offset += rows.Size
		if offset >= f.Length {
			break
		}
		fmt.Println()
	}
	return nil
}

func showFilesInfo(files []*csvfile.CSVFile) {
	for _, file := range files {
		showFileInfo(file)
	}
}

func showFileInfo(file *csvfile.CSVFile) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("Path: ")
	buf.WriteString(file.Path)
	buf.WriteString("\n")
	buf.WriteString("Line/Column delimiters: ")
	buf.WriteString(strconv.Quote(file.Delimiter.String()))
	buf.WriteString("\n")
	buf.WriteString("Length: ")
	buf.WriteString(strconv.FormatInt(file.Length, 10))
	buf.WriteString("\n")
	if len(file.ColumnNames) == 0 {
		buf.WriteString("<no column names found>\n")
	} else {
		buf.WriteString("Column names:\n")
		for _, columnName := range file.ColumnNames {
			buf.WriteString("\t")
			buf.WriteString(strings.Join(columnName.ColumnNames, file.Delimiter.ColumnDelimiter))
			buf.WriteString("\n\t")
			buf.WriteString("Offset: ")
			buf.WriteString(strconv.FormatInt(columnName.Offset, 10))
			buf.WriteString("\n")
		}
	}
	fmt.Println(buf.String())
}

func parseFiles(patterns []string) ([]*csvfile.CSVFile, error) {
	var files []*csvfile.CSVFile
	var errs []error
	for _, pattern := range patterns {
		fileNames, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}

		for _, fileName := range fileNames {
			fz, err := csvfile.ParseCSvFile(fileName)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to parse %q  %v", fileName, err))
				continue
			}
			files = append(files, fz)
		}
	}
	return files, errors.Join(errs...)
}

func exitError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
