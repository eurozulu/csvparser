package datasets

import (
	"github.com/eurozulu/csvparser/utils"
	"testing"
)

func TestRowIterator_ColumnNames(t *testing.T) {
	f, err := ParseCSvFile(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	sets, err := f.DataSets()
	if err != nil {
		t.Fatal(err)
	}
	if len(sets) != 2 {
		t.Fatal("expected 2 data sets")
	}
	ri, err := newRowIterator(sets[0])
	if err != nil {
		t.Fatal(err)
	}
	expect := []string{"one", "two", "three", "four"}
	if !utils.Equals(ri.ColumnNames(), expect) {
		t.Errorf("unexpected column names, expected %v, got %v", expect, ri.ColumnNames())
	}

	ri, err = newRowIterator(sets[1])
	if err != nil {
		t.Fatal(err)
	}
	expect = []string{"id", "first-name", "surname", "department", "email"}
	if !utils.Equals(ri.ColumnNames(), expect) {
		t.Errorf("unexpected column names, expected %v, got %v", expect, ri.ColumnNames())
	}

	ri, err = newRowIterator(sets[1], "id", "email")
	if err != nil {
		t.Fatal(err)
	}
	expect = []string{"id", "email"}
	if !utils.Equals(ri.ColumnNames(), expect) {
		t.Errorf("unexpected column names, expected %v, got %v", expect, ri.ColumnNames())
	}

	ri, err = newRowIterator(sets[1], "id", "unknown-column")
	if err == nil {
		t.Errorf("expected error in second data set for unknown column name")
	}
}

func TestRowIterator_NextRows(t *testing.T) {
	f, err := ParseCSvFile(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	sets, err := f.DataSets()
	if err != nil {
		t.Fatal(err)
	}
	if len(sets) != 2 {
		t.Fatal("expected 2 data sets")
	}
	ri, err := newRowIterator(sets[1])
	if err != nil {
		t.Fatal(err)
	}
	rows, err := ri.NextRows(1)
	if err != nil {
		t.Errorf("error in next rows: %v", err)
	}
	if len(rows) != 1 {
		t.Errorf("expected 1 row, got %d", len(rows))
	}
	expect := []string{"1", "mary", "berry", "Catering", "mb@acme.com"}
	if !utils.Equals(rows[0], expect) {
		t.Errorf("unexpected row, expected %v, got %v", expect, rows[0])
	}
	rows, err = ri.NextRows(999)
	if err != nil {
		t.Errorf("error in next rows: %v", err)
	}
	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}

}
