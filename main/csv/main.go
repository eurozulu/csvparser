package main

import (
	"fmt"
	"github.com/eurozulu/csvparser/csvfile"
	"os"
	"path/filepath"
	"time"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("Took: %v\n", time.Since(start))
	}()

	args := Arguments(os.Args[1:])
	if args.ContainsFlag("help", "?") {
		showhelp()
		return
	}

	patterns := args.Parameters()
	if len(patterns) == 0 {
		showhelp()
		exitError(fmt.Errorf("must provide a filename of the csv file to parse"))
	}
	if flg, ok := args.Flag("line-delimit", "l"); ok {
		if flg == "" {
			exitError(fmt.Errorf("must provide a line delimiter"))
		}
		csvfile.CustomLineDelimiter = flg
	}
	if flg, ok := args.Flag("column-delimit", "c"); ok {
		if flg == "" {
			exitError(fmt.Errorf("must provide a column delimiter"))
		}
		csvfile.CustomColumnDelimiter = flg
	}

	filez, err := parseFiles(patterns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		if len(filez) == 0 {
			return
		}
	}

	if _, ok := args.Flag("info", "i"); ok {
		fileShow{}.ShowFilesInfo(filez)
		return
	}

	colNames, _ := args.Flag("column-names", "n")
	rShow := &rowShow{
		ShowHeaders: args.ContainsFlag("show-headers", "h"),
		ColumnNames: colNames,
	}

	if err = rShow.ShowFiles(filez); err != nil {
		exitError(err)
	}
}

func showhelp() {
	fmt.Println("Usage: csv [options] <csvfile file patter> [<csv file patter>...]")
	fmt.Println("delimitOptions:")
	fmt.Println("  -i, --info\t Display only meta data about the files")
	fmt.Println("  -h, --show-headers\t Display column header names")
	fmt.Println("  -n, --column-names\t specify comma delimited list of column names to display")

	fmt.Println("  -l, --line-delimit\t specify a line delimiter (defaults to new line")
	fmt.Println("  -c, --column-delimit\t specify a line delimiter (defaults to new line")
}

func parseFiles(patterns []string) ([]*csvfile.CSVFile, error) {
	var files []*csvfile.CSVFile
	for _, pattern := range patterns {
		fileNames, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		if len(fileNames) == 0 {
			return nil, fmt.Errorf("no files found with %q", pattern)
		}
		filez, err := csvfile.ParseCSVFiles(fileNames)
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
