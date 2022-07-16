package queryparser

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type QueryParser interface {
	QueryParse(string) error
}

func Parse(r *http.Request, dest any) (err error) {
	v := reflect.ValueOf(dest)
	q := r.URL.Query()
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a pointer to not nil struct")
	}
	v = v.Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fVal := v.Field(i)
		fType := t.Field(i)
		fName := strings.ToLower(fType.Name)
		if queryTag := fType.Tag.Get("query"); queryTag != "" {
			fName = queryTag
		}
		err = parse(q.Get(fName), fVal)
		if err != nil {
			return fmt.Errorf("parse %w", err)
		}
	}
	return nil
}

func parse(stringVal string, destVal reflect.Value) (err error) {
	if stringVal == "" {
		return nil
	}
	if !destVal.CanSet() {
		return fmt.Errorf("field unexported or cannot set value")
	}
	k := destVal.Kind()
	switch {
	case k == reflect.String:
		err = parseString(stringVal, destVal)
	case k >= reflect.Int && k <= reflect.Int64:
		err = parseInt(stringVal, destVal)
	case k >= reflect.Float32 && k <= reflect.Float64:
		err = parseFloat(stringVal, destVal)
	case k >= reflect.Uint && k <= reflect.Uint64:
		err = parseUint(stringVal, destVal)
	case k == reflect.Bool:
		err = parseBool(stringVal, destVal)
	case k == reflect.Slice:
		err = parseSlice(stringVal, destVal)
	default:
		err = parseDefault(stringVal, destVal)
	}
	return err
}

func parseString(in string, dest reflect.Value) error {
	dest.SetString(in)
	return nil
}

func parseInt(in string, dest reflect.Value) error {
	intVal, err := strconv.ParseInt(in, 10, 0)
	if err != nil {
		return fmt.Errorf("parseInt %w", err)
	}
	dest.SetInt(intVal)
	return nil
}

func parseFloat(in string, dest reflect.Value) error {
	floatVal, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return fmt.Errorf("parseFloat %w", err)
	}
	dest.SetFloat(floatVal)
	return nil
}

func parseUint(in string, dest reflect.Value) error {
	uintVal, err := strconv.ParseUint(in, 10, 0)
	if err != nil {
		return fmt.Errorf("parseUint %w", err)
	}
	dest.SetUint(uintVal)
	return nil
}

func parseBool(in string, dest reflect.Value) error {
	boolVal, err := strconv.ParseBool(in)
	if err != nil {
		return fmt.Errorf("parseBool %w", err)
	}
	dest.SetBool(boolVal)
	return nil
}

func parseSlice(in string, dest reflect.Value) error {
	parts := strings.Split(in, ",")
	sliceType := dest.Type().Elem()
	sliceLen := len(parts)
	sliceVal := reflect.MakeSlice(reflect.SliceOf(sliceType), sliceLen, sliceLen)
	for i := 0; i < sliceLen; i++ {
		err := parse(parts[i], sliceVal.Index(i))
		if err != nil {
			return fmt.Errorf("parseSlice %w", err)
		}
	}
	dest.Set(sliceVal)
	return nil
}

func parseDefault(in string, dest reflect.Value) error {
	if dest.Kind() != reflect.Ptr {
		dest = dest.Addr()
	} else if dest.IsNil() {
		dest.Set(reflect.New(dest.Type().Elem()))
	}
	if queryParser, ok := dest.Interface().(QueryParser); ok {
		return queryParser.QueryParse(in)
	}
	return fmt.Errorf("type not supported: %s", dest.Type().Kind())
}
