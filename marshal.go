package bytecodec

import (
	"fmt"
	"reflect"

	"github.com/shimmeringbee/bytecodec/bitbuffer"
)

func Marshal(v interface{}) ([]byte, error) {
	bb := bitbuffer.NewBitBuffer()

	if err := MarshalToBitBuffer(bb, v); err != nil {
		return []byte{}, err
	}

	return bb.Bytes(), nil
}

func MarshalToBitBuffer(bb *bitbuffer.BitBuffer, v interface{}) error {
	val := reflect.Indirect(reflect.ValueOf(v))

	ctx := Context{
		Root:         val,
		CurrentIndex: 0,
	}

	return marshalValue(bb, ctx, "root", val, val, val, "")
}

func marshalValue(bb *bitbuffer.BitBuffer, ctx Context, name string, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) (err error) {
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
		err = marshalArrayOrSlice(bb, ctx, value, root, parent, tags)
	case reflect.String:
		err = marshalString(bb, value, tags)
	case reflect.Ptr:
		err = marshalPtr(bb, ctx, value)
	default:
		err = fmt.Errorf("%w: field '%s' of type '%v'", ErrUnsupportedType, name, kind)
	}

	return
}

func marshalPtr(bb *bitbuffer.BitBuffer, ctx Context, value reflect.Value) error {
	marshaler := reflect.TypeOf((*Marshaler)(nil)).Elem()

	if value.Type().Implements(marshaler) {
		retVals := value.MethodByName("Marshal").Call([]reflect.Value{reflect.ValueOf(bb), reflect.ValueOf(ctx)})

		if retVals[0].IsNil() {
			return nil
		}

		return retVals[0].Interface().(error)
	}

	return fmt.Errorf("%w: field does not support the Marshaler interface", ErrUnsupportedType)
}

func marshalStruct(bb *bitbuffer.BitBuffer, structValue reflect.Value, root reflect.Value) error {
	structType := structValue.Type()

	ctx := Context{
		Root:         structValue,
		CurrentIndex: 0,
	}

	for i := 0; i < structValue.NumField(); i++ {
		value := structValue.Field(i)
		field := structType.Field(i)
		tags := field.Tag
		name := field.Name

		ctx.CurrentIndex = i

		if err := marshalValue(bb, ctx, name, value, root, structValue, tags); err != nil {
			return err
		}
	}

	return nil
}

func marshalArrayOrSlice(bb *bitbuffer.BitBuffer, ctx Context, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) error {
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
		if err := marshalValue(bb, ctx, name, value.Index(i), root, parent, tags); err != nil {
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
	}

	return bb.WriteStringLengthPrefixed(stringValue, stringTag.Endian, int(stringTag.Size))
}

func marshalBool(bb *bitbuffer.BitBuffer, bitSize int, value bool) error {
	byteValue := 0x00

	if value {
		byteValue = 0x01
	}

	return bb.WriteBits(byte(byteValue), bitSize)
}
