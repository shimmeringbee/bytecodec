package bytecodec

import (
	"errors"
	"fmt"
	"github.com/shimmeringbee/bytecodec/bitbuffer"
	"io"
	"math"
	"reflect"
)

func Unmarshal(data []byte, v interface{}) (err error) {
	bb := bitbuffer.NewBitBufferFromBytes(data)
	return UnmarshalFromBitBuffer(bb, v)
}

func UnmarshalFromBitBuffer(bb *bitbuffer.BitBuffer, v interface{}) (err error) {
	val := reflect.Indirect(reflect.ValueOf(v))

	if !val.CanSet() {
		return fmt.Errorf("cannot unmarshall to non pointer")
	}

	ctx := Context{
		Root:         val,
		CurrentIndex: 0,
	}

	if err = unmarshalValue(bb, ctx, "root", val, val, val, ""); err != nil {
		return
	}

	return
}

func unmarshalValue(bb *bitbuffer.BitBuffer, ctx Context, name string, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) (err error) {
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
		err = unmarshalArray(bb, ctx, value, root, parent, tags)
	case reflect.Slice:
		err = unmarshalSlice(bb, ctx, value, root, parent, tags)
	case reflect.Ptr:
		err = unmarshalPtr(bb, ctx, value)
	case reflect.String:
		err = unmarshalString(bb, value, tags)
	default:
		err = fmt.Errorf("%w: field '%s' of type '%v'", UnsupportedType, name, kind)
	}

	return
}

func unmarshalPtr(bb *bitbuffer.BitBuffer, ctx Context, value reflect.Value) error {
	unmarshaler := reflect.TypeOf((*Unmarshaler)(nil)).Elem()

	if value.Type().Implements(unmarshaler) {
		if value.IsNil() {
			e := reflect.New(value.Type().Elem())
			if value.CanSet() {
				value.Set(e)
			}
		}

		retVals := value.MethodByName("Unmarshal").Call([]reflect.Value{reflect.ValueOf(bb), reflect.ValueOf(ctx)})

		if retVals[0].IsNil() {
			return nil
		}

		return retVals[0].Interface().(error)
	} else {
		return fmt.Errorf("%w: field does not support the Marshaler interface", UnsupportedType)
	}

	return nil
}

func unmarshalStruct(bb *bitbuffer.BitBuffer, structValue reflect.Value, root reflect.Value) error {
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

		if err := unmarshalValue(bb, ctx, name, value, root, structValue, tags); err != nil {
			return err
		}
	}

	return nil
}

func unmarshalBool(bb *bitbuffer.BitBuffer, endian bitbuffer.Endian, bitSize int, value reflect.Value) error {
	readValue, err := bb.ReadUint(endian, bitSize)

	if err != nil {
		return err
	}

	value.SetBool(readValue > 0)

	return nil
}

func unmarshalUint(bb *bitbuffer.BitBuffer, endian bitbuffer.Endian, bitSize int, value reflect.Value) error {
	readValue, err := bb.ReadUint(endian, bitSize)

	if err != nil {
		return err
	}

	value.SetUint(readValue)
	return nil
}

func unmarshalArray(bb *bitbuffer.BitBuffer, ctx Context, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) error {
	arraySize, err := readArraySliceLength(bb, tags, value.Len())
	if err != nil {
		return err
	}

	for i := 0; i < arraySize; i++ {
		name := fmt.Sprintf("array[%d]", i)
		if err := unmarshalValue(bb, ctx, name, value.Index(i), root, parent, tags); err != nil {
			return err
		}
	}

	return nil
}

func unmarshalSlice(bb *bitbuffer.BitBuffer, ctx Context, value reflect.Value, root reflect.Value, parent reflect.Value, tags reflect.StructTag) error {
	sliceSize, err := readArraySliceLength(bb, tags, math.MaxInt64)
	if err != nil {
		return err
	}

	value.Set(reflect.MakeSlice(value.Type(), 0, 0))

	for i := 0; i < sliceSize; i++ {
		sliceValue := reflect.New(value.Type().Elem()).Elem()

		name := fmt.Sprintf("slice[%d]", i)
		if err := unmarshalValue(bb, ctx, name, sliceValue, root, parent, tags); err != nil {
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
		readSize, err := bb.ReadUint(length.Endian, int(length.Size))
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

	if stringTag.Termination == Null {
		if str, err := bb.ReadStringNullTerminated(int(stringTag.Size)); err != nil {
			return err
		} else {
			value.SetString(str)
		}
	} else {
		if str, err := bb.ReadStringLengthPrefixed(stringTag.Endian, int(stringTag.Size)); err != nil {
			return err
		} else {
			value.SetString(str)
		}
	}

	return nil
}
