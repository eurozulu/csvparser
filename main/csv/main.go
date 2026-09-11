package main

import (
	"bytes"
	"csvparser/csvfile"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := Arguments(os.Args[1:])
	if args.ContainsFlag("help", "?") {
		showhelp()
		return
	}

	fileNames := args.Parameters()
	if len(fileNames) == 0 {
		showhelp()
		exitError(fmt.Errorf("must provide a filename of the csv file to parse"))
	}
	if flg, ok := args.Flag("line-delimit", "dl"); ok {
		if flg == "" {
			exitError(fmt.Errorf("must provide a line delimiter"))
		}
		csvfile.CustomLineDelimiter = flg
	}
	if flg, ok := args.Flag("column-delimit", "dc"); ok {
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

	if _, ok := args.Flag("file-info", "info"); ok {
		showFilesInfo(filez)
		return
	}
	showColHead := args.ContainsFlag("show-headers", "h")

	if err = showFilesData(filez, showColHead); err != nil {
		exitError(err)
	}
}

func showhelp() {
	fmt.Println("Usage: csv [options] <csvfile file patter> [<csv file patter>...]")
	fmt.Println("Options:")
	fmt.Println("  -info, --file-info\t Display only meta data about the files")
	fmt.Println("  -h, --show-headers\t Display column header names")

	fmt.Println("  -dl, --line-delimit\t specify a line delimiter (defaults to new line")
	fmt.Println("  -cl, --column-delimit\t specify a line delimiter (defaults to new line")
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
	if len(f.ColumnNames) == 0 {
		rows, err := f.RowSet(0)
		if err != nil {
			return err
		}
		if err = showRows(rows, f.Delimiter.ColumnDelimiter); err != nil {
			return err
		}
	}

	for _, header := range f.ColumnNames {
		if showHead {
			fmt.Println(strings.Join(header.ColumnNames, f.Delimiter.ColumnDelimiter))
		}
		rows, err := f.RowSet(header.Offset)
		if err != nil {
			return err
		}
		if err = showRows(rows, f.Delimiter.ColumnDelimiter); err != nil {
			return err
		}
	}
	return nil
}

func showRows(rows *csvfile.RowSet, delimiter string) error {
	for rows.HasRows() {
		r, err := rows.Rows(255)
		if err != nil && err != io.EOF {
			return err
		}
		for _, l := range r {
			fmt.Println(strings.Join(l, delimiter))
		}
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
	buf.WriteString("length: ")
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
			buf.WriteString("offset: ")
			buf.WriteString(strconv.FormatInt(columnName.Offset, 10))
			buf.WriteString("\n")
		}
	}
	fmt.Println(buf.String())
}

func parseFiles(patterns []string) ([]*csvfile.CSVFile, error) {
	var files []*csvfile.CSVFile
	for _, pattern := range patterns {
		filez, err := csvfile.ParseCSVFiles(pattern)
		if err != nil {
			return nil, err
		}
		files = append(files, filez...)
	}
	return files, nil
}

func exitError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
