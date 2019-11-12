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

	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i)
		field := valType.Field(i)
		kind := field.Type.Kind()
		tags := field.Tag

		endianness := tagEndianness(tags)

		switch kind {
		case reflect.Uint8:
			bb.WriteByte(uint8(value.Uint()))
		case reflect.Uint16:
			writeUint(&bb, endianness, 2, value.Uint())
		case reflect.Uint32:
			writeUint(&bb, endianness, 4, value.Uint())
		case reflect.Uint64:
			writeUint(&bb, endianness, 8, value.Uint())

		default:
			return nil, fmt.Errorf("%w: field '%s' of type '%v'", UnsupportedType, val.Type().Field(i).Name, kind)
		}
	}

	return bb.Bytes(), nil
}

func writeUint(bb *bytes.Buffer, endian Endian, size uint8, value uint64) {
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
