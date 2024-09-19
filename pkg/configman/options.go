package configman

import (
	"io"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

type RunnableOption func(cfgDest interface{}) error

func WithEnv() RunnableOption {
	return func(cfgDest interface{}) error {
		structValue := reflect.ValueOf(cfgDest).Elem()
		return reflectFill(structValue, fillByTagStrategy(func(tags reflect.StructTag) (filler string, skip bool) {
			tag, ok := tags.Lookup("env")
			if !ok {
				return "", true
			}

			valueFromEnv, ok := os.LookupEnv(tag)
			return valueFromEnv, !ok
		}))
	}
}

func WithYAMLFile(path string) RunnableOption {
	return func(cfgDest interface{}) error {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		return WithYAML(file)(cfgDest)
	}
}

func WithYAML(content io.Reader) RunnableOption {
	return func(cfgDest interface{}) error {
		return yaml.NewDecoder(content).Decode(cfgDest)
	}
}
