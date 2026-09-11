package encode

import (
	"errors"
	"fmt"
	"github.com/eurozulu/csvparser"
	"io"
	"reflect"
	"strings"
)

type CSVUnmarshaller interface {
	UnmarshalCSV(data []byte) error
}

func UnmarshallCSVSlice[S any](in io.Reader, withHeader bool, delimiters ...*csvparser.Delimiters) ([]S, error) {
	delimit := csvparser.DelimiterOrDefault(delimiters...)
	p := csvparser.NewCsvParser(in, delimit)
	var result []S
	if withHeader {
		if !p.Scan() {
			return result, errors.New("csvparser: no header")
		}
		var s S
		names := fieldNamesOfStruct(s)
		row := p.Row()
		if strings.Join(names, "") != strings.Join(row, "") {
			return result, fmt.Errorf("header mismatch, expected %q but found %q", names, row)
		}
	}

	for p.Scan() {
		row := strings.Join(p.Row(), delimit.ColumnDelimiter)
		var s S
		err := UnmarshallCSV([]byte(row), s, delimit)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, p.Err()
}

func UnmarshallCSV(data []byte, v any, delimiters ...*csvparser.Delimiters) error {
	delimit := csvparser.DelimiterOrDefault(delimiters...)
	if cv, ok := v.(CSVUnmarshaller); ok {
		return cv.UnmarshalCSV(data)
	}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val = reflect.New(reflect.TypeOf(v).Elem())
		}
		val = val.Elem()
	}
	t := val.Type()
	if t.Kind() != reflect.Struct {
		return errors.New("UnmarshallCSV: not a struct")
	}

	row := csvparser.IgnoreQuotedSplit(string(data), delimit.ColumnDelimiter)
	if len(row) == 0 {
		return fmt.Errorf("UnmarshallCSV: empty row")
	}

	//if len(row) > t.NumField() {
	//	return errors.New("UnmarshallCSV: too many cells for fields")
	//}

	var count int
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		if !fld.IsExported() {
			continue
		}
		if count >= len(row) {
			return fmt.Errorf("UnmarshallCSV: not enough cells for field %s", fld.Name)
		}
		cv, err := StringAsType(row[count], fld.Type)
		if err != nil {
			return fmt.Errorf("could not convert %q into type %s for Field %s  %v", row[count], fld.Type.String(), fld.Name, err)
		}
		val.Field(i).Set(cv)
	}
	return nil
}
