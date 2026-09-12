package main

import (
	"fmt"
	"github.com/eurozulu/csvparser/csvfile"
	"github.com/eurozulu/csvparser/utils"
	"io"
	"strings"
)

const rowBlockSize = 255

type rowShow struct {
	ShowHeaders bool
	ColumnNames string
}

func (rs rowShow) ShowFiles(files []*csvfile.CSVFile) error {
	for _, f := range files {
		fmt.Println(f.Path)
		if err := rs.ShowFile(f); err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}

func (rs rowShow) ShowFile(f *csvfile.CSVFile) error {
	if len(f.ColumnHeaders) == 0 {
		rows, err := f.RowSet()
		if err != nil {
			return err
		}
		if err = rs.showRows(rows, f.Delimiter.ColumnDelimiter); err != nil {
			return err
		}
	}

	for _, header := range f.ColumnHeaders {
		names := header.ColumnNames
		if rs.ColumnNames != "" {
			names = strings.Split(rs.ColumnNames, f.Delimiter.ColumnDelimiter)
			if !utils.ContainsAll(header.ColumnNames, names) {
				continue
			}
		}
		if rs.ShowHeaders {
			fmt.Println(strings.Join(names, f.Delimiter.ColumnDelimiter))
		}
		rows, err := f.RowSet(names...)
		if err != nil {
			return err
		}
		if err = rs.showRows(rows, f.Delimiter.ColumnDelimiter); err != nil {
			return err
		}
	}
	return nil
}

func (rs rowShow) showRows(rows *csvfile.RowSet, delimiter string) error {
	for rows.HasRows() {
		r, err := rows.Rows(rowBlockSize)
		if err != nil && err != io.EOF {
			return err
		}
		for _, l := range r {
			fmt.Println(strings.Join(l, delimiter))
		}
	}
	return nil
}
