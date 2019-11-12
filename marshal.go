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
			writeUint16(&bb, endianness, value.Uint())

		default:
			return nil, fmt.Errorf("%w: field '%s' of type '%v'", UnsupportedType, val.Type().Field(i).Name, kind)
		}
	}

	return bb.Bytes(), nil
}

func writeUint16(bb *bytes.Buffer, endian Endian, value uint64) {
	switch endian {
	case BigEndian:
		bb.WriteByte(byte(value >> 8))
		bb.WriteByte(byte(value))
	case LittleEndian:
		bb.WriteByte(byte(value))
		bb.WriteByte(byte(value >> 8))
	}
}
