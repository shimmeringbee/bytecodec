package bytecodec

import (
	"errors"
	"fmt"
	"github.com/shimmeringbee/bytecodec/bitbuffer"
	"io"
	"math"
	"reflect"
	"strings"
)

func Unmarshal(data []byte, v interface{}) (err error) {
	bb := bitbuffer.NewBitBufferFromBytes(data)
	return UnmarshalFromBitBuffer(bb, v)
}

func UnmarshalFromBitBuffer(bb *bitbuffer.BitBuffer, v interface {}) (err error) {
	val := reflect.Indirect(reflect.ValueOf(v))

	if !val.CanSet() {
		return fmt.Errorf("cannot unmarshall to non pointer")
	}

	if err = unmarshalValue(bb, "root", val, val, val, ""); err != nil {
		return
	}

	return
}

func unmarshalValue(bb *bitbuffer.BitBuffer, name string, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) (err error) {
	kind := value.Kind()

	endian := tagEndianness(tags)

	if skip, err := shouldIgnore(tags, root, parent); skip || err != nil {
		return err
	}

	fieldWidth, err := tagFieldWidth(tags)

	if err != nil {
		return
	}

	switch kind {
	case reflect.Bool:
		err = unmarshalBool(bb, endian, fieldWidth.Width(8), value)
	case reflect.Uint8:
		err = unmarshalUint(bb, endian, fieldWidth.Width(8), value)
	case reflect.Uint16:
		err = unmarshalUint(bb, endian, fieldWidth.Width(16), value)
	case reflect.Uint32:
		err = unmarshalUint(bb, endian, fieldWidth.Width(32), value)
	case reflect.Uint64:
		err = unmarshalUint(bb, endian, fieldWidth.Width(64), value)
	case reflect.Struct:
		err = unmarshalStruct(bb, value, root)
	case reflect.Array:
		err = unmarshalArray(bb, value, root, parent, tags)
	case reflect.Slice:
		err = unmarshalSlice(bb, value, root, parent, tags)
	case reflect.Ptr:
		err = unmarshalPtr(bb, value)
	case reflect.String:
		err = unmarshalString(bb, value, tags)
	default:
		err = fmt.Errorf("%w: field '%s' of type '%v'", UnsupportedType, name, kind)
	}

	return
}

func unmarshalPtr(bb *bitbuffer.BitBuffer, value reflect.Value) error {
	unmarshaler := reflect.TypeOf((*Unmarshaler)(nil)).Elem()

	if value.Type().Implements(unmarshaler) {
		if value.IsNil() {
			e := reflect.New(value.Type().Elem())
			if value.CanSet() {
				value.Set(e)
			}
		}

		value.MethodByName("Unmarshal").Call([]reflect.Value{reflect.ValueOf(bb)})
	} else {
		return fmt.Errorf("%w: field does not support the Marshaler interface", UnsupportedType)
	}

	return nil
}

func unmarshalStruct(bb *bitbuffer.BitBuffer, structValue reflect.Value, root reflect.Value) error {
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

func unmarshalBool(bb *bitbuffer.BitBuffer, endian EndianTag, bitSize int, value reflect.Value) error {
	readValue, err := readUintFromBuffer(bb, endian, bitSize)

	if err != nil {
		return err
	}

	value.SetBool(readValue > 0)

	return nil
}

func unmarshalUint(bb *bitbuffer.BitBuffer, endian EndianTag, size int, value reflect.Value) error {
	readValue, err := readUintFromBuffer(bb, endian, size)

	if err != nil {
		return err
	}

	value.SetUint(readValue)
	return nil
}

func readUintFromBuffer(bb *bitbuffer.BitBuffer, endian EndianTag, bitSize int) (uint64, error) {
	if bitSize < 8 {
		data, err := bb.ReadBits(bitSize)
		return uint64(data), err
	} else {
		size := (bitSize + 7) / 8

		if (size * 8) != bitSize {
			return 0, fmt.Errorf("unable to handle arbitary bit widths > 8 bits, %d requested", bitSize)
		}

		readValue := uint64(0)

		switch endian {
		case BigEndian:
			for i := 0; i < size; i++ {
				readByte, err := bb.ReadByte()
				if err != nil {
					return 0, err
				}

				shiftOffset := (size - i - 1) * 8
				readValue |= uint64(readByte) << shiftOffset
			}
		case LittleEndian:
			for i := 0; i < size; i++ {
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
}

func unmarshalArray(bb *bitbuffer.BitBuffer, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) error {
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

func unmarshalSlice(bb *bitbuffer.BitBuffer, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) error {
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

func readArraySliceLength(bb *bitbuffer.BitBuffer, tags reflect.StructTag, max int) (int, error) {
	length, err := tagSlicePrefix(tags)
	if err != nil {
		return 0, err
	}

	if length.HasPrefix() {
		readSize, err := readUintFromBuffer(bb, length.Endian, int(length.Size))
		if err != nil {
			return 0, err
		}

		return int(readSize), nil
	} else {
		return max, nil
	}
}

func unmarshalString(bb *bitbuffer.BitBuffer, value reflect.Value, tags reflect.StructTag) error {
	stringTag, err := tagStringType(tags)
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
		stringLength, err := readUintFromBuffer(bb, stringTag.Endian, int(stringTag.Size))
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
