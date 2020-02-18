package bytecodec

import (
	"fmt"
	"reflect"
	"strconv"
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

		i, err := findValue(includeBase, includeIf.FieldPath)

		if err != nil {
			return false, err
		}

		switch v := i.(type) {
		case bool:
			stringValue := includeIf.Value

			if stringValue == "" {
				stringValue = "true"
			}

			tagVal, err := strconv.ParseBool(stringValue)

			switch includeIf.Operation {
			case Equal:
				return tagVal != v, err
			case NotEqual:
				return tagVal == v, err
			default:
				return false, fmt.Errorf("includeIf path could not be parsed: unable to compare end parameter (unknown comparison for bool)")
			}
		case uint8:
			return compareUint(uint64(v), includeIf)
		case uint16:
			return compareUint(uint64(v), includeIf)
		case uint32:
			return compareUint(uint64(v), includeIf)
		case uint64:
			return compareUint(v, includeIf)
		default:
			return false, fmt.Errorf("includeIf path could not be parsed: unable to compare end parameter (unknown type)")
		}
	}

	return false, nil
}

func compareUint(v uint64, includeIf IncludeIfTag) (bool, error) {
	stringValue := includeIf.Value

	if stringValue == "" {
		stringValue = "0"
	}

	tagVal, err := strconv.ParseUint(stringValue, 10, 64)

	switch includeIf.Operation {
	case Equal:
		return tagVal != v, err
	case NotEqual:
		return tagVal == v, err
	default:
		return false, fmt.Errorf("includeIf path could not be parsed: unable to compare end parameter (unknown comparison for bool)")
	}
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

			return value.Interface(), nil
		}
	}

	return false, fmt.Errorf("includeIf path could not be parsed: %s not found", thisLevel)
}
