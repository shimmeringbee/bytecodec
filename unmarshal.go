package bytecodec

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
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
	case reflect.Array:
		err = unmarshallArray(bb, value, tags)
	case reflect.Slice:
		err = unmarshallSlice(bb, value, tags)
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
	readValue, err := readUintFromBuffer(bb, endian, size)

	if err != nil {
		return err
	}

	value.SetUint(readValue)
	return nil
}

func readUintFromBuffer(bb *bytes.Buffer, endian EndianTag, size uint8) (uint64, error) {
	readValue := uint64(0)

	switch endian {
	case BigEndian:
		for i := uint8(0); i < size; i++ {
			readByte, err := bb.ReadByte()
			if err != nil {
				return 0, err
			}

			shiftOffset := (size - i - 1) * 8
			readValue |= uint64(readByte) << shiftOffset
		}
	case LittleEndian:
		for i := uint8(0); i < size; i++ {
			readByte, err := bb.ReadByte()
			if err != nil {
				return 0, err
			}

			shiftOffset := i * 8
			readValue |= uint64(readByte) << shiftOffset
		}
	}

	return readValue, nil
}

func unmarshallArray(bb *bytes.Buffer, value reflect.Value, tags reflect.StructTag) error {
	arraySize, err := readArraySliceLength(bb, tags, value.Len())
	if err != nil {
		return err
	}

	for i := 0; i < arraySize; i++ {
		name := fmt.Sprintf("array[%d]", i)
		if err := unmarshalValue(bb, name, value.Index(i), tags); err != nil {
			return err
		}
	}

	return nil
}

func unmarshallSlice(bb *bytes.Buffer, value reflect.Value, tags reflect.StructTag) error {
	sliceSize, err := readArraySliceLength(bb, tags, math.MaxInt64)
	if err != nil {
		return err
	}

	value.Set(reflect.MakeSlice(value.Type(), 0, 0))

	for i := 0; i < sliceSize; i++ {
		sliceValue := reflect.New(value.Type().Elem()).Elem()

		name := fmt.Sprintf("slice[%d]", i)
		if err := unmarshalValue(bb, name, sliceValue, tags); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}

		value.Set(reflect.Append(value, sliceValue))
	}

	return nil
}

func readArraySliceLength(bb *bytes.Buffer, tags reflect.StructTag, max int) (int, error) {
	length, err := tagLength(tags)
	if err != nil {
		return 0, err
	}

	if length.HasLength() {
		readSize, err := readUintFromBuffer(bb, length.Endian, length.Size)
		if err != nil {
			return 0, err
		}

		return int(readSize), nil
	} else {
		return max, nil
	}
}
