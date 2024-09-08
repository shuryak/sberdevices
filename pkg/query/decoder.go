package query

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

type Decoder struct {
	values url.Values
}

func NewDecoder(v url.Values) *Decoder {
	return &Decoder{values: v}
}

func (d *Decoder) Decode(dest interface{}) error {
	lastQuery := make(map[string]string)
	for key := range d.values {
		lastQuery[key] = d.values[key][len(d.values[key])-1]
	}

	typ := reflect.TypeOf(dest)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	structValue := reflect.ValueOf(dest).Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := structValue.Field(i)
		tag := field.Tag.Get("query")

		queryValue, ok := lastQuery[tag]
		if !ok || tag == "" || !value.CanSet() {
			continue
		}

		switch value.Kind() {
		case reflect.Int:
			intValue, _ := strconv.Atoi(queryValue)
			value.SetInt(int64(intValue))
		case reflect.String:
			value.SetString(queryValue)
		default:
			return fmt.Errorf("unsupported query type: %s", tag)
		}
	}

	return nil
}
