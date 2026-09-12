package datasets

import (
	"errors"
	"fmt"
	"github.com/eurozulu/csvparser/utils"
	"testing"
)

const testFilePath = "../.testdata/1_mytest.csv"
const testFileLength = 263

func TestParseCSvFile(t *testing.T) {
	f, err := ParseCSvFile(testFilePath)
	if err != nil {
		t.Fatalf("Error parsing CSV file: %v", err)
	}
	if f.Path != testFilePath {
		t.Errorf("unexpected CSV file path. expected %s, got %s", testFilePath, f.Path)
	}
	expectLen := int64(testFileLength)
	if f.Length != expectLen {
		t.Errorf("unexpected CSV file length. expected %d, got %d", expectLen, f.Length)
	}
	expectDelimit := "\n ::"
	gotDelimit := f.Delimiter.String()
	if expectDelimit != gotDelimit {
		t.Errorf("unexpected Delimiter. expected %s, got %s", expectDelimit, gotDelimit)
	}
}

func TestCSVFile_DataSets(t *testing.T) {
	f, err := ParseCSvFile(testFilePath)
	if err != nil {
		t.Fatalf("Error parsing CSV file: %v", err)
	}
	sets, err := f.DataSets()
	if err != nil {
		t.Fatalf("Error parsing CSV file datasets: %v", err)
	}
	if len(sets) != 2 {
		t.Errorf("Unexpected number of set data. expected %d, got %d", 2, len(sets))
	}

	if err := checkDataSet(sets[0],
		"\n ::", 0, 106, "one", "two", "three", "four"); err != nil {
		t.Errorf("error in first data set: %v", err)
	}
	if err := checkDataSet(sets[1],
		"\n ::", 107, 156, "id", "first-name", "surname", "department", "email"); err != nil {
		t.Errorf("error in second data set: %v", err)
	}
}

func checkDataSet(set *DataSet, expectDelimit string, expectOffset, expectLength int64, expectColumn ...string) error {
	var errs []error
	if set.Delimiter.String() != expectDelimit {
		errs = append(errs, fmt.Errorf("unexpected Delimiter. expected %s, got %s", expectDelimit, set.Delimiter.String()))
	}
	if set.Offset != expectOffset {
		errs = append(errs, fmt.Errorf("unexpected Offset. expected %d, got %d", expectOffset, set.Offset))
	}
	if set.Size != expectLength {
		errs = append(errs, fmt.Errorf("unexpected Length. expected %d, got %d", expectLength, set.Size))
	}
	if !utils.Equals(set.Columns, expectColumn) {
		errs = append(errs, fmt.Errorf("unexpected Columns. expected %v, got %v", expectColumn, set.Columns))
	}
	return errors.Join(errs...)
}
