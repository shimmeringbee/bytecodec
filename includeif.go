package bytecodec

import (
	"fmt"
	"reflect"
)

func shouldIgnore(tags reflect.StructTag, root reflect.Value, parent reflect.Value) (bool, error) {
	includeIf, err := tagIncludeIf(tags)

	if err != nil {
		return false, err
	}

	if includeIf.HasIncludeIf() {
		includeBase := root

		if includeIf.Relative {
			includeBase = parent
		}

		val, err := findValue(includeBase, includeIf.FieldPath)

		if err != nil {
			return false, err
		}

		switch val {
		case val.(bool):
			return includeIf.Value != val, err
		}
	}

	return false, nil
}

func findValue(structValue reflect.Value, path []string) (interface{}, error) {
	structType := structValue.Type()

	thisLevel := path[0]
	remainingPath := path[1:]

	for i := 0; i < structValue.NumField(); i++ {
		value := structValue.Field(i)
		field := structType.Field(i)
		name := field.Name

		if name == thisLevel {
			if len(remainingPath) >= 1 {
				if value.Kind() != reflect.Struct {
					return false, fmt.Errorf("includeIf path could not be parsed: %s is not a struct", name)
				}

				return findValue(value, remainingPath)
			}

			if value.Kind() != reflect.Bool {
				return false, fmt.Errorf("includeIf path could not be parsed: %s is not a boolean (end parameter)", name)
			}

			return value.Interface(), nil
		}
	}

	return false, fmt.Errorf("includeIf path could not be parsed: %s not found", thisLevel)
}
