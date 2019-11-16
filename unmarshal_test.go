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

	t.Run("verify nested struct is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8
			Two struct {
				Three uint8
			}
		}

		expectedStruct := &StructUnderTest{One: 0x01, Two: struct{ Three uint8 }{Three: 0x03}}
		data, _ := Marshall(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshall(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify nested struct with unmarshalable type errors", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8
			Two struct {
				Three chan bool
			}
		}

		data := []byte{0x00, 0x01}

		actualStruct := StructUnderTest{}
		err := Unmarshall(data, &actualStruct)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, UnsupportedType))
		assert.Equal(t, "unsupported type: field 'Three' of type 'chan'", err.Error())
	})

	t.Run("verify the unmarshal of non struct is little endian", func(t *testing.T) {
		expectedStruct := uint32(0x80010203)

		data, _ := Marshall(expectedStruct)

		actualStruct := uint32(0)
		err := Unmarshall(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify a slice of bytes unmarshals", func(t *testing.T) {
		type StructUnderTest struct {
			One []byte
		}

		expectedStruct := &StructUnderTest{One: []byte{0x55, 0xaa}}
		data, _ := Marshall(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshall(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify an array of bytes unmarshals", func(t *testing.T) {
		type StructUnderTest struct {
			One [2]byte
			Two byte
		}

		expectedStruct := &StructUnderTest{One: [2]byte{0x55, 0xaa}, Two: 0x02}
		data, _ := Marshall(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshall(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify an array with unmarshalable type errors", func(t *testing.T) {
		type StructUnderTest struct {
			One [2]chan bool
		}

		instance := &StructUnderTest{}
		var data []byte
		err := Unmarshall(data, instance)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, UnsupportedType))
		assert.Equal(t, "unsupported type: field 'array[0]' of type 'chan'", err.Error())
	})

	t.Run("verify a slice of uint16s obeys big endian annotation", func(t *testing.T) {
		type StructUnderTest struct {
			One []uint16 `bcendian:"big"`
		}

		expectedStruct := &StructUnderTest{One: []uint16{0x8001}}
		data, _ := Marshall(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshall(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify custom type definition unmarshals", func(t *testing.T) {
		type NetworkAddress [8]byte

		type StructUnderTest struct {
			One NetworkAddress
		}

		expectedStruct := &StructUnderTest{One: NetworkAddress{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}}
		data, _ := Marshall(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshall(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})
}
