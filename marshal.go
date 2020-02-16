package bytecodec

import (
	"errors"
	"fmt"
	"github.com/shimmeringbee/bytecodec/bitbuffer"
	"math"
	"reflect"
)

var UnsupportedType = errors.New("unsupported type")

func Marshal(v interface{}) ([]byte, error) {
	bb := bitbuffer.NewBitBuffer()

	val := reflect.Indirect(reflect.ValueOf(v))

	if err := marshalValue(bb, "root", val, val, val, ""); err != nil {
		return nil, err
	}

	return bb.Bytes(), nil
}

func marshalValue(bb *bitbuffer.BitBuffer, name string, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) (err error) {
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
		err = marshalBool(bb, fieldWidth.Width(8), value.Bool())
	case reflect.Uint8:
		err = marshalUint(bb, endian, fieldWidth.Width(8), value.Uint())
	case reflect.Uint16:
		err = marshalUint(bb, endian, fieldWidth.Width(16), value.Uint())
	case reflect.Uint32:
		err = marshalUint(bb, endian, fieldWidth.Width(32), value.Uint())
	case reflect.Uint64:
		err = marshalUint(bb, endian, fieldWidth.Width(64), value.Uint())
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

func marshalStruct(bb *bitbuffer.BitBuffer, structValue reflect.Value, root reflect.Value) error {
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

func marshalArrayOrSlice(bb *bitbuffer.BitBuffer, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) error {
	length, err := tagSlicePrefix(tags)
	if err != nil {
		return err
	}

	if length.HasPrefix() {
		if err := marshalUint(bb, length.Endian, int(length.Size), uint64(value.Len())); err != nil {
			return err
		}
	}

	for i := 0; i < value.Len(); i++ {
		name := fmt.Sprintf("array[%d]", i)
		if err := marshalValue(bb, name, value.Index(i), root, parent, tags); err != nil {
			return err
		}
	}

	return nil
}

func marshalString(bb *bitbuffer.BitBuffer, value reflect.Value, tags reflect.StructTag) error {
	stringTag, err := tagStringType(tags)
	if err != nil {
		return err
	}

	stringBytes := []byte(value.String())
	stringLength := len(stringBytes)

	if stringTag.Termination == Null {
		if stringTag.Size > 0 && (stringLength+1) > int(stringTag.Size) {
			return fmt.Errorf("string too large to fit in padding allocated")
		}

		if err := bb.WriteString(value.String()); err != nil {
			return err
		}

		if err := bb.WriteByte(0); err != nil {
			return err
		}

		for i := 0; i < int(stringTag.Size)-(stringLength+1); i++ {
			if err := bb.WriteByte(0); err != nil {
				return err
			}
		}
	} else {
		maxSize := int(math.Pow(2, float64(stringTag.Size)))

		if stringLength > maxSize {
			return fmt.Errorf("string too large to be represented by prefixed length")
		}

		if err := marshalUint(bb, stringTag.Endian, int(stringTag.Size), uint64(stringLength)); err != nil {
			return err
		}

		if err := bb.WriteString(value.String()); err != nil {
			return err
		}
	}

	return nil
}

func marshalBool(bb *bitbuffer.BitBuffer, bitSize int, value bool) error {
	byteValue := 0x00

	if value {
		byteValue = 0x01
	}

	return bb.WriteBits(byte(byteValue), bitSize)
}

func marshalUint(bb *bitbuffer.BitBuffer, endian EndianTag, bitSize int, value uint64) error {
	maxValue := math.Pow(2, float64(bitSize)) - 1

	if float64(value) > maxValue {
		return fmt.Errorf("cannot marshal value %d into %d bit field", value, bitSize)
	}

	if bitSize < 8 {
		if err := bb.WriteBits(byte(value), bitSize); err != nil {
			return err
		}
	} else {
		size := (bitSize + 7) / 8

		if (size * 8) != bitSize {
			return fmt.Errorf("unable to handle arbitary bit widths > 8 bits, %d requested", bitSize)
		}

		switch endian {
		case BigEndian:
			for i := 0; i < size; i++ {
				shiftOffset := (size - i - 1) * 8

				if err := bb.WriteByte(byte(value >> shiftOffset)); err != nil {
					return err
				}
			}
		case LittleEndian:
			for i := 0; i < size; i++ {
				shiftOffset := i * 8
				if err := bb.WriteByte(byte(value >> shiftOffset)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
