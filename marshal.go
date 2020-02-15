package bytecodec

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"reflect"
)

var UnsupportedType = errors.New("unsupported type")

func Marshal(v interface{}) ([]byte, error) {
	bb := bytes.Buffer{}

	val := reflect.Indirect(reflect.ValueOf(v))

	if err := marshalValue(&bb, "root", val, val, val, ""); err != nil {
		return nil, err
	}

	return bb.Bytes(), nil
}

func marshalValue(bb *bytes.Buffer, name string, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) (err error) {
	kind := value.Kind()

	endian := tagEndianness(tags)

	if skip, err := shouldIgnore(tags, root, parent); skip || err != nil {
		return err
	}

	switch kind {
	case reflect.Bool:
		marshalBool(bb, endian, value.Bool())
	case reflect.Uint8:
		marshalUint(bb, endian, 1, value.Uint())
	case reflect.Uint16:
		marshalUint(bb, endian, 2, value.Uint())
	case reflect.Uint32:
		marshalUint(bb, endian, 4, value.Uint())
	case reflect.Uint64:
		marshalUint(bb, endian, 8, value.Uint())
	case reflect.Struct:
		err = marshalStruct(bb, value, root)
	case reflect.Array, reflect.Slice:
		err = marshalArrayOrSlice(bb, value, root, parent, tags)
	case reflect.String:
		err = marshalString(bb, value, tags)
	default:
		err = fmt.Errorf("%w: field '%s' of type '%v'", UnsupportedType, name, kind)
	}

	return
}

func marshalStruct(bb *bytes.Buffer, structValue reflect.Value, root reflect.Value) error {
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		value := structValue.Field(i)
		field := structType.Field(i)
		tags := field.Tag
		name := field.Name

		if err := marshalValue(bb, name, value, root, structValue, tags); err != nil {
			return err
		}
	}

	return nil
}

func marshalArrayOrSlice(bb *bytes.Buffer, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) error {
	length, err := tagLength(tags)
	if err != nil {
		return err
	}

	if length.HasLength() {
		marshalUint(bb, length.Endian, length.Size, uint64(value.Len()))
	}

	for i := 0; i < value.Len(); i++ {
		name := fmt.Sprintf("array[%d]", i)
		if err := marshalValue(bb, name, value.Index(i), root, parent, tags); err != nil {
			return err
		}
	}

	return nil
}

func marshalString(bb *bytes.Buffer, value reflect.Value, tags reflect.StructTag) error {
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

		marshalUint(bb, stringTag.Endian, (stringTag.Size+7)/8, uint64(stringLength))
		bb.WriteString(value.String())
	}

	return nil
}

func marshalBool(bb *bytes.Buffer, endian EndianTag, value bool) {
	byteValue := 0x00

	if value {
		byteValue = 0x01
	}

	bb.WriteByte(byte(byteValue))
}

func marshalUint(bb *bytes.Buffer, endian EndianTag, size uint8, value uint64) {
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
