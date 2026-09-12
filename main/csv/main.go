package main

import (
	"fmt"
	"github.com/eurozulu/csvparser"
	"github.com/eurozulu/csvparser/datasets"
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

	var delimit *csvparser.Delimiters
	if flg, ok := args.Flag("line-delimit", "l"); ok {
		if flg == "" {
			exitError(fmt.Errorf("must provide a line delimiter"))
		}
		d := csvparser.DelimiterOrDefault()
		d.LineDelimiter = flg
		delimit = &d
	}
	if flg, ok := args.Flag("column-delimit", "c"); ok {
		if flg == "" {
			exitError(fmt.Errorf("must provide a column delimiter"))
		}
		if delimit == nil {
			d := csvparser.DelimiterOrDefault()
			delimit = &d
		}
		delimit.ColumnDelimiter = flg
	}

	var dlz []csvparser.Delimiters
	if delimit != nil {
		dlz = append(dlz, *delimit)
	}
	filez, err := parseFiles(patterns, dlz...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		if len(filez) == 0 {
			return
		}
	}

	if _, ok := args.Flag("info", "i"); ok {
		if err := (&fileShow{}).ShowFilesInfo(filez); err != nil {
			exitError(err)
		}
		return
	}

	colNames, _ := args.Flag("column-names", "n")
	rShow := &rowShow{
		ShowHeaders: args.ContainsFlag("show-headers", "h"),
		ColumnNames: csvparser.IgnoreQuotedSplit(colNames, ","),
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

func parseFiles(patterns []string, delimiter ...csvparser.Delimiters) ([]*datasets.CSVFile, error) {
	var files []*datasets.CSVFile
	for _, pattern := range patterns {
		fileNames, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		if len(fileNames) == 0 {
			return nil, fmt.Errorf("no files found with %q", pattern)
		}
		filez, err := datasets.ParseCSVFiles(fileNames, delimiter...)
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
