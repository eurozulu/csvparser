package encode

import (
	"encoding"
	"fmt"
	"github.com/eurozulu/csvparser"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

var (
	textUnmarshalerType   = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	textMarshalerType     = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	binaryUnmarshalerType = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()
	binaryMarshalerType   = reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()
	urlType               = reflect.TypeOf((*url.URL)(nil)).Elem()
)

func TypeAsString(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Implements(textMarshalerType) {
		b, err := v.(encoding.TextMarshaler).MarshalText()
		if err != nil {
			panic(err)
		}
		return string(b)
	}
	if t.Implements(binaryMarshalerType) {
		b, err := v.(encoding.BinaryMarshaler).MarshalBinary()
		if err != nil {
			panic(err)
		}
		return string(b)
	}

	switch t.Kind() {
	case reflect.String:
		return csvparser.doubleQuote(v.(string))
	case reflect.Bool:
		return strconv.FormatBool(v.(bool))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.(int64), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.(uint64), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.(float64), 'g', -1, 64)
	case reflect.Slice:
		return csvparser.doubleQuote(sliceAsString(v))

	default:
		return fmt.Sprintf("%v", v)
	}
}

func StringAsType(s string, tp reflect.Type) (reflect.Value, error) {
	t := tp
	if t.Kind() == reflect.Ptr {
		v, err := StringAsType(s, t.Elem())
		if err != nil {
			return v, err
		}
		return v.Addr(), nil
	}

	switch t.Kind() {
	case reflect.String:
		return reflect.ValueOf(s), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(b), nil
	case reflect.Int:
		i, err := strconv.Atoi(s)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Int8:
		i, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Int16:
		i, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Int32:
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Uint, reflect.Uint32:
		i, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Uint8:
		i, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Uint16:
		i, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Uint64:
		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Float32:
		i, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Float64:
		i, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Slice:
		return stringAsSliceType(s, t)
	default:
		if tp.Implements(textUnmarshalerType) {
			v := reflect.New(tp)
			if err := v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(s)); err != nil {
				return reflect.Value{}, fmt.Errorf("Failed to unmarshal %q as text %v", s, err)
			}
			return v.Elem(), nil
		}
		if tp.Implements(binaryUnmarshalerType) {
			v := reflect.New(tp)
			if err := v.Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary([]byte(s)); err != nil {
				return reflect.Value{}, fmt.Errorf("Failed to unmarshal %q as binary %v", s, err)
			}
			return v.Elem(), nil
		}
		if urlType.AssignableTo(tp) {
			u, err := url.Parse(s)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("Failed to parse %q as URL %v", s, err)
			}
			return reflect.ValueOf(u).Elem(), nil
		}
		return reflect.Value{}, fmt.Errorf("Unsupported type: %s", t.String())
	}
}

func sliceAsString(v any) string {
	val := reflect.ValueOf(v)
	l := val.Len()
	values := make([]string, l)
	for i := 0; i < l; i++ {
		values[i] = TypeAsString(val.Index(i).Interface())
	}
	return strings.Join(values, ",")
}

func stringAsSliceType(s string, tp reflect.Type) (reflect.Value, error) {
	ss := strings.Split(s, ",")
	sv := reflect.MakeSlice(tp, len(ss), len(ss))
	t := tp.Elem()
	for i, sz := range ss {
		v, err := StringAsType(sz, t)
		if err != nil {
			return sv, err
		}
		sv.Index(i).Set(v)
	}
	return sv, nil
}
