package encode

import (
	"reflect"
	"errors"
	"bytes"
	"strings"
	"io"
	"csvparser"
)

type CSVMarshaller interface {
	MarshalCSV() ([]byte, error)
}

func fieldNamesOfStruct(v any) []string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("csvparser.ColumnNames: not a struct")
	}
	names := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		if !fld.IsExported() {
			continue
		}
		if tag, ok := fld.Tag.Lookup("csv"); ok {
			if tag == "-" {
				continue
			}
			names = append(names, tag)
		} else {
			names = append(names, strings.ToLower(fld.Name))
		}
	}
	return names
}

func MarshallCSVSlice[S any](out io.Writer, includeHeaders bool, v []S, delimiters ...*csvparser.Delimiters) error {
	delimit := csvparser.DelimiterOrDefault(delimiters...)

	if includeHeaders {
		var s S
		if cols := fieldNamesOfStruct(s); len(cols) > 0 {
			if _, err := out.Write([]byte(strings.Join(cols, delimit.ColumnDelimiter))); err != nil {
				return err
			}
			if _, err := out.Write([]byte(delimit.LineDelimiter)); err != nil {
				return err
			}
		}
	}
	for _, vv := range v {
		data, err := MarshallCSV(vv, delimit)
		if err != nil {
			return err
		}
		if _, err = out.Write(data); err != nil {
			return err
		}
		if _, err = out.Write([]byte(delimit.LineDelimiter)); err != nil {
			return err
		}
	}
	return nil
}

func MarshallCSV(v any, delimiters ...*csvparser.Delimiters) ([]byte, error) {
	delimit := csvparser.DelimiterOrDefault(delimiters...)

	if cv, ok := v.(CSVMarshaller); ok {
		return cv.MarshalCSV()
	}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	t := val.Type()
	if t.Kind() != reflect.Struct {
		return nil, errors.New("csvparser.ColumnNames: not a struct")
	}

	buf := bytes.NewBuffer(nil)
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		if !fld.IsExported() {
			continue
		}
		tags := strings.Split(fld.Tag.Get("csv"), ",")
		if tags[0] == "-" {
			continue
		}
		if len(tags) > 1 && strings.EqualFold(tags[1], "omitempty") {
			if val.Field(i).IsZero() {
				continue
			}
		}

		sv := TypeAsString(val.Field(i).Interface())
		if buf.Len() > 0 {
			buf.WriteString(delimit.ColumnDelimiter)
		}
		buf.WriteString(sv)
	}
	if buf.Len() > 0 {
		buf.WriteString(delimit.LineDelimiter)
	}
	return buf.Bytes(), nil
}
