package bytecodec

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

var UnsupportedType = errors.New("unsupported type")

func Marshall(v interface{}) ([]byte, error) {
	bb := bytes.Buffer{}

	val := reflect.Indirect(reflect.ValueOf(v))

	if err := marshalValue(&bb, "root", val, LittleEndian); err != nil {
		return nil, err
	}

	return bb.Bytes(), nil
}

func marshalStruct(bb *bytes.Buffer, value reflect.Value) error {
	structType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		value := value.Field(i)
		field := structType.Field(i)
		tags := field.Tag
		name := field.Name

		endianness := tagEndianness(tags)

		if err := marshalValue(bb, name, value, endianness); err != nil {
			return err
		}
	}

	return nil
}

func marshalValue(bb *bytes.Buffer, name string, value reflect.Value, endian Endian) (err error) {
	kind := value.Kind()

	switch kind {
	case reflect.Uint8:
		marshallUint(bb, endian, 1, value.Uint())
	case reflect.Uint16:
		marshallUint(bb, endian, 2, value.Uint())
	case reflect.Uint32:
		marshallUint(bb, endian, 4, value.Uint())
	case reflect.Uint64:
		marshallUint(bb, endian, 8, value.Uint())
	case reflect.Struct:
		err = marshalStruct(bb, value)
	case reflect.Array, reflect.Slice:
		err = marshallArrayOrSlice(bb, value)
	default:
		err = fmt.Errorf("%w: field '%s' of type '%v'", UnsupportedType, name, kind)
	}

	return
}

func marshallArrayOrSlice(bb *bytes.Buffer, value reflect.Value) error {
	for i := 0; i < value.Len(); i++ {
		name := fmt.Sprintf("array[%d]", i)
		if err := marshalValue(bb, name, value.Index(i), LittleEndian); err != nil {
			return err
		}
	}

	return nil
}

func marshallUint(bb *bytes.Buffer, endian Endian, size uint8, value uint64) {
	switch endian {
	case BigEndian:
		for i := uint8(0); i < size; i++ {
			shiftOffset := (size - i - 1) * 8
			bb.WriteByte(byte(value >> shiftOffset))
		}
	case LittleEndian:
		for i := uint8(0); i < size; i++ {
			shiftOffset := i * 8
			bb.WriteByte(byte(value >> shiftOffset))
		}
	}
}
