package encode

import (
	"errors"
	"fmt"
	"github.com/eurozulu/csvparser"
	"io"
	"reflect"
)

type CSVUnmarshaller interface {
	UnmarshalCSV(data []byte) error
}

func UnmarshallCSVSlice[S any](in io.Reader, delimiters ...csvparser.Delimiters) ([]S, error) {
	delimit := csvparser.DelimiterOrDefault(delimiters...)
	scn := csvparser.NewIgnoreQuotedScanner(in, delimit.LineDelimiter)
	var result []S
	for scn.Scan() {
		var s S
		err := UnmarshallCSV(scn.Bytes(), &s, delimit)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, scn.Err()
}

func UnmarshallCSV(data []byte, v any, delimiters ...csvparser.Delimiters) error {
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
			return fmt.Errorf("could not convert %q into type %s for Field %s  %v",
				row[count], fld.Type.String(), fld.Name, err)
		}
		val.Field(i).Set(cv)
	}
	return nil
}
