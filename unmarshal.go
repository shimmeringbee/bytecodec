package bytecodec

import (
	"bytes"
	"fmt"
	"reflect"
)

func Unmarshall(data []byte, v interface{}) (err error) {
	bb := bytes.NewBuffer(data)

	val := reflect.Indirect(reflect.ValueOf(v))

	if err = unmarshalValue(bb, v, "root", val, ""); err != nil {
		return
	}

	return
}

func unmarshalValue(bb *bytes.Buffer, v interface{}, name string, value reflect.Value, tags reflect.StructTag) (err error) {
	kind := value.Kind()

	switch kind {
	case reflect.Struct:
		err = unmarshalStruct(bb, v, value)
	default:
		err = fmt.Errorf("%w: field '%s' of type '%v'", UnsupportedType, name, kind)
	}

	return
}

func unmarshalStruct(bb *bytes.Buffer, v interface{}, value reflect.Value) error {
	structType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		value := value.Field(i)
		field := structType.Field(i)
		tags := field.Tag
		name := field.Name

		if err := unmarshalValue(bb, v, name, value, tags); err != nil {
			return err
		}
	}

	return nil
}
