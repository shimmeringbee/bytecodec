package bytecodec

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
)

func Unmarshal(data []byte, v interface{}) (err error) {
	bb := bytes.NewBuffer(data)

	val := reflect.Indirect(reflect.ValueOf(v))

	if !val.CanSet() {
		return fmt.Errorf("cannot unmarshall to non pointer")
	}

	if err = unmarshalValue(bb, "root", val, val, val, ""); err != nil {
		return
	}

	return
}

func unmarshalValue(bb *bytes.Buffer, name string, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) (err error) {
	kind := value.Kind()

	endian := tagEndianness(tags)

	if skip, err := shouldIgnore(tags, root, parent); skip || err != nil {
		return err
	}

	switch kind {
	case reflect.Bool:
		err = unmarshalBool(bb, endian, value)
	case reflect.Uint8:
		err = unmarshalUint(bb, endian, 1, value)
	case reflect.Uint16:
		err = unmarshalUint(bb, endian, 2, value)
	case reflect.Uint32:
		err = unmarshalUint(bb, endian, 4, value)
	case reflect.Uint64:
		err = unmarshalUint(bb, endian, 8, value)
	case reflect.Struct:
		err = unmarshalStruct(bb, value, root)
	case reflect.Array:
		err = unmarshalArray(bb, value, root, parent, tags)
	case reflect.Slice:
		err = unmarshalSlice(bb, value, root, parent, tags)
	case reflect.String:
		err = unmarshalString(bb, value, tags)
	default:
		err = fmt.Errorf("%w: field '%s' of type '%v'", UnsupportedType, name, kind)
	}

	return
}

func unmarshalStruct(bb *bytes.Buffer, structValue reflect.Value, root reflect.Value) error {
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		value := structValue.Field(i)
		field := structType.Field(i)
		tags := field.Tag
		name := field.Name

		if err := unmarshalValue(bb, name, value, root, structValue, tags); err != nil {
			return err
		}
	}

	return nil
}

func unmarshalBool(bb *bytes.Buffer, endian EndianTag, value reflect.Value) error {
	readValue, err := readUintFromBuffer(bb, endian, 1)

	if err != nil {
		return err
	}

	value.SetBool(readValue > 0)

	return nil
}

func unmarshalUint(bb *bytes.Buffer, endian EndianTag, size uint8, value reflect.Value) error {
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

func unmarshalArray(bb *bytes.Buffer, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) error {
	arraySize, err := readArraySliceLength(bb, tags, value.Len())
	if err != nil {
		return err
	}

	for i := 0; i < arraySize; i++ {
		name := fmt.Sprintf("array[%d]", i)
		if err := unmarshalValue(bb, name, value.Index(i), root, parent, tags); err != nil {
			return err
		}
	}

	return nil
}

func unmarshalSlice(bb *bytes.Buffer, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) error {
	sliceSize, err := readArraySliceLength(bb, tags, math.MaxInt64)
	if err != nil {
		return err
	}

	value.Set(reflect.MakeSlice(value.Type(), 0, 0))

	for i := 0; i < sliceSize; i++ {
		sliceValue := reflect.New(value.Type().Elem()).Elem()

		name := fmt.Sprintf("slice[%d]", i)
		if err := unmarshalValue(bb, name, sliceValue, root, parent, tags); err != nil {
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

func unmarshalString(bb *bytes.Buffer, value reflect.Value, tags reflect.StructTag) error {
	stringTag, err := tagString(tags)
	if err != nil {
		return err
	}

	sb := strings.Builder{}

	if stringTag.Termination == Null {
		maxLength := math.MaxInt64

		if stringTag.Size != 0 {
			maxLength = int(stringTag.Size)
		}

		i := 0

		readByte := ^byte(0)

		for ; i < maxLength; i++ {
			readByte, err = bb.ReadByte()
			if err != nil {
				return err
			}

			if readByte == 0 {
				i++
				break
			}

			sb.WriteByte(readByte)
		}

		if readByte != 0 {
			return errors.New("no null termination found in string")
		}

		for ; i < int(stringTag.Size); i++ {
			_, err := bb.ReadByte()
			if err != nil {
				return err
			}
		}

	} else {
		stringLength, err := readUintFromBuffer(bb, stringTag.Endian, (stringTag.Size+7)/8)
		if err != nil {
			return err
		}

		for i := 0; i < int(stringLength); i++ {
			readByte, err := bb.ReadByte()
			if err != nil {
				return err
			}

			sb.WriteByte(readByte)
		}
	}

	value.SetString(sb.String())
	return nil
}
