package configman

import (
	"fmt"
	"reflect"
	"strconv"
)

func fillFromTagStrategy(tagName string) fillStrategy {
	return func(structField reflect.StructField, structValue reflect.Value) (filler string, skip bool, err error) {
		filler, ok := structField.Tag.Lookup(tagName)
		return filler, !ok, nil
	}
}

func fillByTagStrategy(tagValueHandler tagValueHandler) fillStrategy {
	return func(structField reflect.StructField, structValue reflect.Value) (filler string, skip bool, err error) {
		filler, skip = tagValueHandler(structField.Tag)
		return filler, skip, nil
	}
}

type tagValueHandler func(tags reflect.StructTag) (filler string, skip bool)

func reflectFill(structElem reflect.Value, strategy fillStrategy) error {
	typ := structElem.Type()

	for i := 0; i < structElem.NumField(); i++ {
		filler, skip, err := strategy(typ.Field(i), structElem.Field(i))
		if err != nil {
			return err
		}

		field := structElem.Field(i)
		kind := field.Kind()

		if kind == reflect.Struct {
			return reflectFill(field, strategy)
		}

		if skip {
			continue
		}

		switch kind {
		case reflect.String:
			field.SetString(filler)
		case reflect.Bool:
			if filler == "YES" || filler == "Yes" || filler == "yes" {
				field.SetBool(true)
				break
			}
			v, _ := strconv.ParseBool(filler)
			field.SetBool(v)
		case reflect.Int:
			v, _ := strconv.Atoi(filler)
			field.SetInt(int64(v))
		default:
			return fmt.Errorf("unsupported type of env %s: %s", typ, kind)
		}
	}

	return nil
}

type fillStrategy func(
	structField reflect.StructField,
	structValue reflect.Value,
) (filler string, skip bool, err error)
