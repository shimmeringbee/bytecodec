package bytecodec

import (
	"bytes"
	"fmt"
	"reflect"
)

func Unmarshall(data []byte, v interface{}) (err error) {
	bb := bytes.NewBuffer(data)

	val := reflect.Indirect(reflect.ValueOf(v))

	if !val.CanSet() {
		return fmt.Errorf("cannot unmarshall to non pointer")
	}

	if err = unmarshalValue(bb, "root", val, ""); err != nil {
		return
	}

	return
}

func unmarshalValue(bb *bytes.Buffer, name string, value reflect.Value, tags reflect.StructTag) (err error) {
	kind := value.Kind()

	endian := tagEndianness(tags)

	switch kind {
	case reflect.Uint8:
		err = unmarshallUint(bb, endian, 1, value)
	case reflect.Uint16:
		err = unmarshallUint(bb, endian, 2, value)
	case reflect.Uint32:
		err = unmarshallUint(bb, endian, 4, value)
	case reflect.Uint64:
		err = unmarshallUint(bb, endian, 8, value)
	case reflect.Struct:
		err = unmarshalStruct(bb, value)
	default:
		err = fmt.Errorf("%w: field '%s' of type '%v'", UnsupportedType, name, kind)
	}

	return
}

func unmarshalStruct(bb *bytes.Buffer, value reflect.Value) error {
	structType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		value := value.Field(i)
		field := structType.Field(i)
		tags := field.Tag
		name := field.Name

		if err := unmarshalValue(bb, name, value, tags); err != nil {
			return err
		}
	}

	return nil
}

func unmarshallUint(bb *bytes.Buffer, endian EndianTag, size uint8, value reflect.Value) error {
	readValue := uint64(0)

	switch endian {
	case BigEndian:
		for i := uint8(0); i < size; i++ {
			readByte, err := bb.ReadByte()
			if err != nil {
				return err
			}

			shiftOffset := (size - i - 1) * 8
			readValue |= uint64(readByte) << shiftOffset
		}
	case LittleEndian:
		for i := uint8(0); i < size; i++ {
			readByte, err := bb.ReadByte()
			if err != nil {
				return err
			}

			shiftOffset := i * 8
			readValue |= uint64(readByte) << shiftOffset
		}
	}

	value.SetUint(readValue)
	return nil
}
