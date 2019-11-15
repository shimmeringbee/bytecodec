package bytecodec

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshall(t *testing.T) {
	t.Run("verifies that error is thrown on unsupported type", func(t *testing.T) {
		type StructUnderTest struct {
			One chan bool
		}

		var data []byte
		instance := &StructUnderTest{}
		err := Unmarshall(data, instance)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, UnsupportedType))
		assert.Equal(t, "unsupported type: field 'One' of type 'chan'", err.Error())
	})

	t.Run("verify non pointers raise an error", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8
		}

		actualStruct := StructUnderTest{}
		err := Unmarshall([]byte{0x00}, actualStruct)

		assert.Error(t, err)
	})

	t.Run("verify byte and uint8 unmarshals", func(t *testing.T) {
		type StructUnderTest struct {
			One byte
			Two uint8
		}

		expectedStruct := StructUnderTest{One: 0x55, Two: 0xaa}
		data, _ := Marshall(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshall(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify uint16 LE is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint16 `bcendian:"little"`
		}

		expectedStruct := StructUnderTest{One: 0x8001}
		data, _ := Marshall(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshall(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify uint16 BE is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint16 `bcendian:"big"`
		}

		expectedStruct := StructUnderTest{One: 0x8001}
		data, _ := Marshall(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshall(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify uint32 LE is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint32 `bcendian:"little"`
		}

		expectedStruct := StructUnderTest{One: 0x80010203}
		data, _ := Marshall(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshall(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify uint32 BE is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint32 `bcendian:"big"`
		}

		expectedStruct := StructUnderTest{One: 0x80010203}
		data, _ := Marshall(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshall(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify uint64 LE is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint64 `bcendian:"little"`
		}

		expectedStruct := StructUnderTest{One: 0x8001020304050607}
		data, _ := Marshall(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshall(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify uint64 BE is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint64 `bcendian:"big"`
		}

		expectedStruct := StructUnderTest{One: 0x8001020304050607}
		data, _ := Marshall(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshall(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})
}
