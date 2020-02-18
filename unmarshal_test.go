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
		err := Unmarshal(data, instance)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, UnsupportedType))
		assert.Equal(t, "unsupported type: field 'One' of type 'chan'", err.Error())
	})

	t.Run("verify non pointers raise an error", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8
		}

		actualStruct := StructUnderTest{}
		err := Unmarshal([]byte{0x00}, actualStruct)

		assert.Error(t, err)
	})

	t.Run("verify bool unmarshals", func(t *testing.T) {
		type StructUnderTest struct {
			One bool
		}

		expectedStruct := StructUnderTest{One: true}
		data, _ := Marshal(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshal(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify byte and uint8 unmarshals", func(t *testing.T) {
		type StructUnderTest struct {
			One byte
			Two uint8
		}

		expectedStruct := StructUnderTest{One: 0x55, Two: 0xaa}
		data, _ := Marshal(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshal(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify uint16 LE is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint16 `bcendian:"little"`
		}

		expectedStruct := StructUnderTest{One: 0x8001}
		data, _ := Marshal(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshal(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify uint16 BE is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint16 `bcendian:"big"`
		}

		expectedStruct := StructUnderTest{One: 0x8001}
		data, _ := Marshal(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshal(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify uint32 LE is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint32 `bcendian:"little"`
		}

		expectedStruct := StructUnderTest{One: 0x80010203}
		data, _ := Marshal(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshal(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify uint32 BE is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint32 `bcendian:"big"`
		}

		expectedStruct := StructUnderTest{One: 0x80010203}
		data, _ := Marshal(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshal(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify uint64 LE is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint64 `bcendian:"little"`
		}

		expectedStruct := StructUnderTest{One: 0x8001020304050607}
		data, _ := Marshal(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshal(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify uint64 BE is unmarshalled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint64 `bcendian:"big"`
		}

		expectedStruct := StructUnderTest{One: 0x8001020304050607}
		data, _ := Marshal(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshal(data, &actualStruct)

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
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

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
		err := Unmarshal(data, &actualStruct)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, UnsupportedType))
		assert.Equal(t, "unsupported type: field 'Three' of type 'chan'", err.Error())
	})

	t.Run("verify the unmarshal of non struct is little endian", func(t *testing.T) {
		expectedStruct := uint32(0x80010203)

		data, _ := Marshal(expectedStruct)

		actualStruct := uint32(0)
		err := Unmarshal(data, &actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify a slice of bytes unmarshals", func(t *testing.T) {
		type StructUnderTest struct {
			One []byte
		}

		expectedStruct := &StructUnderTest{One: []byte{0x55, 0xaa}}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify an array of bytes unmarshals", func(t *testing.T) {
		type StructUnderTest struct {
			One [2]byte
			Two byte
		}

		expectedStruct := &StructUnderTest{One: [2]byte{0x55, 0xaa}, Two: 0x02}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify an array with unmarshalable type errors", func(t *testing.T) {
		type StructUnderTest struct {
			One [2]chan bool
		}

		instance := &StructUnderTest{}
		var data []byte
		err := Unmarshal(data, instance)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, UnsupportedType))
		assert.Equal(t, "unsupported type: field 'array[0]' of type 'chan'", err.Error())
	})

	t.Run("verify a slice of uint16s obeys big endian annotation", func(t *testing.T) {
		type StructUnderTest struct {
			One []uint16 `bcendian:"big"`
		}

		expectedStruct := &StructUnderTest{One: []uint16{0x8001}}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify custom type definition unmarshals", func(t *testing.T) {
		type NetworkAddress [8]byte

		type StructUnderTest struct {
			One NetworkAddress
		}

		expectedStruct := &StructUnderTest{One: NetworkAddress{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify slices support implicit length annotations, uint8", func(t *testing.T) {
		type StructUnderTest struct {
			One []byte `bcsliceprefix:"8"`
		}

		expectedStruct := &StructUnderTest{One: []byte{0x55, 0xaa}}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify arrays support implicit length annotations, uint8", func(t *testing.T) {
		type StructUnderTest struct {
			One [2]byte `bcsliceprefix:"8"`
		}

		expectedStruct := &StructUnderTest{One: [2]byte{0x55, 0xaa}}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify slices support implicit length annotations, uint16, big endian", func(t *testing.T) {
		type StructUnderTest struct {
			One []byte `bcsliceprefix:"16,big"`
		}

		expectedStruct := &StructUnderTest{One: []byte{0x55, 0xaa}}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify slices support implicit length annotations, uint16, default little endian", func(t *testing.T) {
		type StructUnderTest struct {
			One []byte `bcsliceprefix:"16"`
		}

		expectedStruct := &StructUnderTest{One: []byte{0x55, 0xaa}}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify size uint8 prefixed string unmarshals", func(t *testing.T) {
		type StructUnderTest struct {
			One string
		}

		expectedStruct := &StructUnderTest{One: "abc"}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify size uint16 big endian prefixed string unmarshals", func(t *testing.T) {
		type StructUnderTest struct {
			One string `bcstringtype:"prefix,16,big"`
		}

		expectedStruct := &StructUnderTest{One: "abc"}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify null terminated string unmarshals", func(t *testing.T) {
		type StructUnderTest struct {
			One string `bcstringtype:"null"`
		}

		expectedStruct := &StructUnderTest{One: "abc"}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify null terminated string with padding unmarshals, skipping padding", func(t *testing.T) {
		type StructUnderTest struct {
			One string `bcstringtype:"null,8"`
			Two byte
		}

		expectedStruct := &StructUnderTest{One: "abc", Two: 0x80}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("verify null terminated string with padding errors if null terminator could not be read", func(t *testing.T) {
		type StructUnderTest struct {
			One string `bcstringtype:"null,4"`
		}

		data := []byte{0x61, 0x62, 0x63, 0x64}

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.Error(t, err)
	})

	t.Run("verifies that unmarshalling struct with includeif is obeyed", func(t *testing.T) {
		type StructUnderTest struct {
			One   bool
			Two   uint8 `bcincludeif:"One"`
			Three uint8
		}

		inputStruct := &StructUnderTest{One: false, Two: 2, Three: 3}
		expectedStruct := &StructUnderTest{One: false, Two: 0, Three: 3}

		data, _ := Marshal(inputStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("unmarshalling two 3 bit uints", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8 `bcfieldwidth:"3"`
			Two uint8 `bcfieldwidth:"3"`
		}

		expectedStruct := &StructUnderTest{One: 0b101, Two: 0b101}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("writing a multibyte integer which is not a standard size", func(t *testing.T) {
		type StructUnderTest struct {
			One uint32 `bcfieldwidth:"24"`
		}

		expectedStruct := &StructUnderTest{One: 0x00aabbcc}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("an error is thrown when a non byte aligned width over 8 bits is requested", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8 `bcfieldwidth:"9"`
		}

		instance := &StructUnderTest{One: 0}
		err := Unmarshal([]byte{}, instance)

		assert.Error(t, err)
	})

	t.Run("two 1 bit bools and a 6 bit uint are read", func(t *testing.T) {
		type StructUnderTest struct {
			One   bool  `bcfieldwidth:"1"`
			Two   bool  `bcfieldwidth:"1"`
			Three uint8 `bcfieldwidth:"6"`
		}

		expectedStruct := &StructUnderTest{One: true, Two: true, Three: 0x2d}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("6 bit uint, array with 2 bit prefix are read", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8  `bcfieldwidth:"6"`
			Two []byte `bcsliceprefix:"2"`
		}

		expectedStruct := &StructUnderTest{One: 0x2d, Two: []byte{0x00, 0x01}}
		data, _ := Marshal(expectedStruct)

		actualStruct := &StructUnderTest{}
		err := Unmarshal(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})

	t.Run("supports pointers which implement Unmarshaller interface", func(t *testing.T) {
		type StructUnderTest struct {
			One *CustomField
		}

		instance := &StructUnderTest{}
		bytes := []byte{0x03, 'O', 'N', 'E'}

		assert.Nil(t, instance.One)

		err := Unmarshal(bytes, instance)

		assert.NoError(t, err)
		assert.NotNil(t, instance.One)
		assert.Equal(t, uint8(1), instance.One.Value)
	})
}
