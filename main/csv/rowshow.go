package main

import (
	"fmt"
	"github.com/eurozulu/csvparser/datasets"
	"github.com/eurozulu/csvparser/utils"
	"io"
	"strings"
)

const rowBlockSize = 255

type rowShow struct {
	ShowHeaders bool
	ColumnNames []string
}

func (rs rowShow) ShowFiles(files []*datasets.CSVFile) error {
	for _, f := range files {
		fmt.Printf("File: %s\n", f.Path)
		if err := rs.ShowFile(f); err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}

func (rs rowShow) ShowFile(f *datasets.CSVFile) error {
	sets, err := f.DataSets()
	if err != nil {
		return err
	}

	for _, set := range sets {
		if len(rs.ColumnNames) > 0 && !utils.ContainsAll(set.Columns, rs.ColumnNames) {
			//set doesn't contain required columns
			continue
		}
		names := rs.ColumnNames
		if len(names) == 0 {
			names = set.Columns
		}
		if rs.ShowHeaders {
			fmt.Println(strings.Join(names, f.Delimiter.ColumnDelimiter))
		}
		rows, err := set.Rows(names...)
		if err != nil {
			return err
		}
		if err = rs.showRows(rows, f.Delimiter.ColumnDelimiter); err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}

func (rs rowShow) showRows(rows *datasets.RowIterator, delimiter string) error {
	for rows.HasNextRows() {
		r, err := rows.NextRows(rowBlockSize)
		if err != nil && err != io.EOF {
			return err
		}
		for _, l := range r {
			fmt.Println(strings.Join(l, delimiter))
		}
	}
	return nil
}
