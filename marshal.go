package bytecodec

import (
	"errors"
	"fmt"
	"github.com/shimmeringbee/bytecodec/bitbuffer"
	"reflect"
)

var UnsupportedType = errors.New("unsupported type")

func Marshal(v interface{}) ([]byte, error) {
	bb := bitbuffer.NewBitBuffer()

	if err := MarshalToBitBuffer(bb, v); err != nil {
		return []byte{}, err
	}

	return bb.Bytes(), nil
}

func MarshalToBitBuffer(bb *bitbuffer.BitBuffer, v interface{}) error {
	val := reflect.Indirect(reflect.ValueOf(v))

	return marshalValue(bb, "root", val, val, val, "")
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
		err = bb.WriteUint(value.Uint(), endian, fieldWidth.Width(8))
	case reflect.Uint16:
		err = bb.WriteUint(value.Uint(), endian, fieldWidth.Width(16))
	case reflect.Uint32:
		err = bb.WriteUint(value.Uint(), endian, fieldWidth.Width(32))
	case reflect.Uint64:
		err = bb.WriteUint(value.Uint(), endian, fieldWidth.Width(64))
	case reflect.Struct:
		err = marshalStruct(bb, value, root)
	case reflect.Array, reflect.Slice:
		err = marshalArrayOrSlice(bb, value, root, parent, tags)
	case reflect.String:
		err = marshalString(bb, value, tags)
	case reflect.Ptr:
		err = marshalPtr(bb, value)
	default:
		err = fmt.Errorf("%w: field '%s' of type '%v'", UnsupportedType, name, kind)
	}

	return
}

func marshalPtr(bb *bitbuffer.BitBuffer, value reflect.Value) error {
	marshaler := reflect.TypeOf((*Marshaler)(nil)).Elem()

	if value.Type().Implements(marshaler) {
		retVals := value.MethodByName("Marshal").Call([]reflect.Value{reflect.ValueOf(bb)})

		if retVals[0].IsNil() {
			return nil
		}

		return retVals[0].Interface().(error)
	} else {
		return fmt.Errorf("%w: field does not support the Marshaler interface", UnsupportedType)
	}
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
		if err := bb.WriteUint(uint64(value.Len()), length.Endian, int(length.Size)); err != nil {
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

	stringValue := value.String()

	if stringTag.Termination == Null {
		return bb.WriteStringNullTerminated(stringValue, int(stringTag.Size))
	} else {
		return bb.WriteStringLengthPrefixed(stringValue, stringTag.Endian, int(stringTag.Size))
	}
}

func marshalBool(bb *bitbuffer.BitBuffer, bitSize int, value bool) error {
	byteValue := 0x00

	if value {
		byteValue = 0x01
	}

	return bb.WriteBits(byte(byteValue), bitSize)
}
