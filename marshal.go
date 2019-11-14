package bytecodec

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"reflect"
)

var UnsupportedType = errors.New("unsupported type")

func Marshall(v interface{}) ([]byte, error) {
	bb := bytes.Buffer{}

	val := reflect.Indirect(reflect.ValueOf(v))

	if err := marshalValue(&bb, "root", val, ""); err != nil {
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

		if err := marshalValue(bb, name, value, tags); err != nil {
			return err
		}
	}

	return nil
}

func marshalValue(bb *bytes.Buffer, name string, value reflect.Value, tags reflect.StructTag) (err error) {
	kind := value.Kind()

	endian := tagEndianness(tags)

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
		err = marshallArrayOrSlice(bb, value, tags)
	case reflect.String:
		err = marshallString(bb, value, tags)
	default:
		err = fmt.Errorf("%w: field '%s' of type '%v'", UnsupportedType, name, kind)
	}

	return
}

func marshallArrayOrSlice(bb *bytes.Buffer, value reflect.Value, tags reflect.StructTag) error {
	length, err := tagLength(tags)
	if err != nil {
		return err
	}

	if length.ShouldInsert() {
		marshallUint(bb, length.Endian, length.Size, uint64(value.Len()))
	}

	for i := 0; i < value.Len(); i++ {
		name := fmt.Sprintf("array[%d]", i)
		if err := marshalValue(bb, name, value.Index(i), tags); err != nil {
			return err
		}
	}

	return nil
}

func marshallString(bb *bytes.Buffer, value reflect.Value, tags reflect.StructTag) error {
	stringTag, err := tagString(tags)
	if err != nil {
		return err
	}

	stringBytes := []byte(value.String())
	stringLength := len(stringBytes)

	if stringTag.Termination == Null {
		if stringTag.Size > 0 && (stringLength+1) > int(stringTag.Size) {
			return fmt.Errorf("string too large to fit in padding allocated")
		}

		bb.WriteString(value.String())
		bb.WriteByte(0)

		for i := 0; i < int(stringTag.Size)-(stringLength+1); i++ {
			bb.WriteByte(0)
		}
	} else {
		maxSize := int(math.Pow(2, float64(stringTag.Size)))

		if stringLength > maxSize {
			return fmt.Errorf("string too large to be represented by prefixed length")
		}

		marshallUint(bb, stringTag.Endian, stringTag.Size, uint64(stringLength))
		bb.WriteString(value.String())
	}

	return nil
}

func marshallUint(bb *bytes.Buffer, endian EndianTag, size uint8, value uint64) {
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
