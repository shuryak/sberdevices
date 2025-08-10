package configman

import (
	"fmt"
	"reflect"
	"strconv"
)

func Collect(cfgDest interface{}, opts ...RunnableOption) error {
	typ := reflect.TypeOf(cfgDest)
	if typ.Kind() != reflect.Ptr {
		return fmt.Errorf("config dest must be a pointer to a struct")
	}

	typ = typ.Elem()

	if len(opts) == 0 {
		return nil
	}

	for _, opt := range opts {
		if err := opt(cfgDest); err != nil {
			return err
		}
	}

	structValue := reflect.ValueOf(cfgDest).Elem()
	return reflectFill(
		structValue,
		func(structField reflect.StructField, structValue reflect.Value) (filler string, skip bool, err error) {
			if structValue.IsZero() {
				if reqFlag, _ := strconv.ParseBool(structField.Tag.Get("required")); reqFlag {
					return "", false, fmt.Errorf("field \"%s\" is required", structField.Name)
				}

				return fillFromTagStrategy("default")(structField, structValue)
			}

			return "", true, nil
		},
	)
}
