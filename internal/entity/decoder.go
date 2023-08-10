package entity

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	ErrNotFound = errors.New("key not found")
)

const (
	_key = "rac"

	_formatDateWoTZ = "2006-01-02T15:04:05"
)

// Helper func for parsing incominf line of text
func GetKeyValue(line string, delimeter rune) (k, v string, err error) {
	cleanLine := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return unicode.ToLower(r)
		}

		return -1
	}, line)

	if pos := strings.IndexRune(cleanLine, delimeter); pos != -1 {
		return strings.Trim(cleanLine[:pos], " "), strings.Trim(cleanLine[pos+1:], " "), nil
	}

	return "", "", ErrNotFound
}

// Converting lines of strings to object
// TODO: use codegen, not reflect - it's slow
func Unmarshal(lines []string, v any) error {
	rt := reflect.TypeOf(v)
	_ = reflect.New(rt)
	vv := reflect.ValueOf(v)

	for _, line := range lines {
		key, value, err := GetKeyValue(line, ':')

		if err != nil && errors.Is(err, ErrNotFound) {
			continue
		}

		if err != nil {
			return err
		}

		for i := 0; i < rt.Elem().NumField(); i++ {
			f := rt.Elem().Field(i)
			val, ok := f.Tag.Lookup(_key)
			if ok && val == key {
				fv := vv.Elem().FieldByName(f.Name)

				switch fv.Interface().(type) {
				case time.Time:
					t, err := time.ParseInLocation(_formatDateWoTZ, strings.ToUpper(value), time.UTC)

					if err == nil {
						fv.Set(reflect.ValueOf(t))
					}
				case int:
					vi, err := strconv.Atoi(value)

					if err == nil {
						fv.SetInt(int64(vi))
					}
				default:
					fv.Set(reflect.ValueOf(value))
				}

				continue
			}
		}
	}

	if vv.Elem().IsZero() {
		return ErrNotFound
	}

	return nil
}
