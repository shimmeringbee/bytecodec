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
	valType := val.Type()

	if err := marshalStruct(&bb, val, valType); err != nil {
		return nil, err
	}

	return bb.Bytes(), nil
}

func marshalStruct(bb *bytes.Buffer, v reflect.Value, t reflect.Type) error {
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		field := t.Field(i)
		kind := field.Type.Kind()
		tags := field.Tag
		name := field.Name

		endianness := tagEndianness(tags)

		if err := marshalValue(bb, name, value, kind, endianness); err != nil {
			return err
		}
	}

	return nil
}

func marshalValue(bb *bytes.Buffer, n string, v reflect.Value, k reflect.Kind, e Endian) error {
	switch k {
	case reflect.Uint8:
		marshallUint(bb, e, 1, v.Uint())
	case reflect.Uint16:
		marshallUint(bb, e, 2, v.Uint())
	case reflect.Uint32:
		marshallUint(bb, e, 4, v.Uint())
	case reflect.Uint64:
		marshallUint(bb, e, 8, v.Uint())

	default:
		return fmt.Errorf("%w: field '%s' of type '%v'", UnsupportedType, n, k)
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
