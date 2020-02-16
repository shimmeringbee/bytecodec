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

		boolVal, err := findBooleanValue(includeBase, includeIf.FieldPath)

		return includeIf.Value != boolVal, err
	}

	return false, nil
}

func findBooleanValue(structValue reflect.Value, path []string) (bool, error) {
	structType := structValue.Type()

	thisLevel := path[0]
	remainingPath := path[1:]

	for i := 0; i < structValue.NumField(); i++ {
		value := structValue.Field(i)
		field := structType.Field(i)
		name := field.Name

		if name == thisLevel {
			if len(remainingPath) > 1 {
				if value.Kind() != reflect.Struct {
					return false, fmt.Errorf("includeIf path could not be parsed: %s is not a struct", name)
				}

				return findBooleanValue(value, remainingPath)
			}

			if value.Kind() != reflect.Bool {
				return false, fmt.Errorf("includeIf path could not be parsed: %s is not a boolean (end parameter)", name)
			}

			return value.Bool(), nil
		}
	}

	return false, fmt.Errorf("includeIf path could not be parsed: %s not found", thisLevel)
}
