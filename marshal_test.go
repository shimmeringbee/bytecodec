package bytecodec

import (
	"errors"
	"github.com/shimmeringbee/bytecodec/bitbuffer"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMarshal(t *testing.T) {
	t.Run("verifies that error is thrown on unsupported type", func(t *testing.T) {
		type StructUnderTest struct {
			One chan bool
		}

		instance := &StructUnderTest{}
		_, err := Marshal(instance)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, UnsupportedType))
		assert.Equal(t, "unsupported type: field 'One' of type 'chan'", err.Error())
	})

	t.Run("verify bool marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One bool
		}

		instance := &StructUnderTest{One: true}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify byte and uint8 marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One byte
			Two uint8
		}

		instance := &StructUnderTest{One: 0x55, Two: 0xaa}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x55, 0xaa}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify uint16 LE is marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One uint16 `bcendian:"little"`
		}

		instance := &StructUnderTest{One: 0x8001}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x01, 0x80}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify uint16 BE is marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One uint16 `bcendian:"big"`
		}

		instance := &StructUnderTest{One: 0x8001}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x80, 0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify uint32 LE is marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One uint32 `bcendian:"little"`
		}

		instance := &StructUnderTest{One: 0x80010203}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x03, 0x02, 0x01, 0x80}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify uint32 BE is marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One uint32 `bcendian:"big"`
		}

		instance := &StructUnderTest{One: 0x80010203}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x80, 0x01, 0x02, 0x03}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify uint64 LE is marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One uint64 `bcendian:"little"`
		}

		instance := &StructUnderTest{One: 0x8001020304050607}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x80}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify uint64 BE is marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One uint64 `bcendian:"big"`
		}

		instance := &StructUnderTest{One: 0x8001020304050607}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x80, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify nested struct is marshaled", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8
			Two struct {
				Three uint8
			}
		}

		instance := &StructUnderTest{One: 0x01, Two: struct{ Three uint8 }{Three: 0x03}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x01, 0x03}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify nested struct with unmarshalable type errors", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8
			Two struct {
				Three chan bool
			}
		}

		instance := &StructUnderTest{One: 0x01, Two: struct{ Three chan bool }{Three: nil}}
		_, err := Marshal(instance)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, UnsupportedType))
		assert.Equal(t, "unsupported type: field 'Three' of type 'chan'", err.Error())
	})

	t.Run("verify the marshal of non struct is little endian", func(t *testing.T) {
		instance := uint32(0x80010203)

		actualBytes, err := Marshal(instance)
		expectedBytes := []byte{0x03, 0x02, 0x01, 0x80}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify a slice of bytes marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One []byte
		}

		instance := &StructUnderTest{One: []byte{0x55, 0xaa}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x55, 0xaa}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify an array of bytes marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One [2]byte
		}

		instance := &StructUnderTest{One: [2]byte{0x55, 0xaa}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x55, 0xaa}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify an array with marshalable type errors", func(t *testing.T) {
		type StructUnderTest struct {
			One [2]chan bool
		}

		instance := &StructUnderTest{}
		_, err := Marshal(instance)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, UnsupportedType))
		assert.Equal(t, "unsupported type: field 'array[0]' of type 'chan'", err.Error())
	})

	t.Run("verify a slice of uint16s obeys big endian annotation", func(t *testing.T) {
		type StructUnderTest struct {
			One []uint16 `bcendian:"big"`
		}

		instance := &StructUnderTest{One: []uint16{0x8001}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x80, 0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify custom type definition marshals", func(t *testing.T) {
		type NetworkAddress [8]byte

		type StructUnderTest struct {
			One NetworkAddress
		}

		instance := &StructUnderTest{One: NetworkAddress{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify slices support implicit length annotations, uint8", func(t *testing.T) {
		type StructUnderTest struct {
			One []byte `bcsliceprefix:"8"`
		}

		instance := &StructUnderTest{One: []byte{0x55, 0xaa}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x02, 0x55, 0xaa}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify arrays support implicit length annotations, uint8", func(t *testing.T) {
		type StructUnderTest struct {
			One [2]byte `bcsliceprefix:"8"`
		}

		instance := &StructUnderTest{One: [2]byte{0x55, 0xaa}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x02, 0x55, 0xaa}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify slices support implicit length annotations, uint16, big endian", func(t *testing.T) {
		type StructUnderTest struct {
			One []byte `bcsliceprefix:"16,big"`
		}

		instance := &StructUnderTest{One: []byte{0x55, 0xaa}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x00, 0x02, 0x55, 0xaa}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify slices support implicit length annotations, uint16, default little endian", func(t *testing.T) {
		type StructUnderTest struct {
			One []byte `bcsliceprefix:"16"`
		}

		instance := &StructUnderTest{One: []byte{0x55, 0xaa}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x02, 0x00, 0x55, 0xaa}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify size uint8 prefixed string marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One string
		}

		instance := &StructUnderTest{One: "abc"}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x03, 0x61, 0x62, 0x63}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify size uint16 big endian prefixed string marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One string `bcstringtype:"prefix,16,big"`
		}

		instance := &StructUnderTest{One: "abc"}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x00, 0x03, 0x61, 0x62, 0x63}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify prefixed string errors if string is longer than prefix can represent", func(t *testing.T) {
		type StructUnderTest struct {
			One string
		}

		instance := &StructUnderTest{One: strings.Repeat("a", 257)}
		_, err := Marshal(instance)

		assert.Error(t, err)
	})

	t.Run("verify null terminated string marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One string `bcstringtype:"null"`
		}

		instance := &StructUnderTest{One: "abc"}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x61, 0x62, 0x63, 0x00}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify null terminated string with padding marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One string `bcstringtype:"null,8"`
		}

		instance := &StructUnderTest{One: "abc"}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x61, 0x62, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify null terminated string with padding errors if no room for null terminator", func(t *testing.T) {
		type StructUnderTest struct {
			One string `bcstringtype:"null,4"`
		}

		instance := &StructUnderTest{One: "abcd"}
		_, err := Marshal(instance)

		assert.Error(t, err)
	})

	t.Run("verify struct with includeIf omits a field if required field is not true", func(t *testing.T) {
		type StructUnderTest struct {
			One bool
			Two uint8 `bcincludeif:".One"`
		}

		instance := &StructUnderTest{One: false, Two: 2}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x00}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify struct with includeIf marshals a field if required field is true", func(t *testing.T) {
		type StructUnderTest struct {
			One bool
			Two uint8 `bcincludeif:".One"`
		}

		instance := &StructUnderTest{One: true, Two: 2}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x01, 0x02}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify struct with includeIf correctly handles absolute references and excludes if false", func(t *testing.T) {
		type Nested struct {
			One bool
			Two uint8 `bcincludeif:".One"`
		}

		type StructUnderTest struct {
			One    bool
			Nested Nested
		}

		instance := &StructUnderTest{One: false, Nested: Nested{One: true, Two: 2}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x00, 0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify struct with includeIf correctly handles absolute references and includes if true", func(t *testing.T) {
		type Nested struct {
			One bool
			Two uint8 `bcincludeif:".One"`
		}

		type StructUnderTest struct {
			One    bool
			Nested Nested
		}

		instance := &StructUnderTest{One: true, Nested: Nested{One: false, Two: 2}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x01, 0x00, 0x02}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify struct with includeIf handles relative reference and excludes if false", func(t *testing.T) {
		type Nested struct {
			One bool
			Two uint8 `bcincludeif:"One"`
		}

		type StructUnderTest struct {
			One    bool
			Nested Nested
		}

		instance := &StructUnderTest{One: true, Nested: Nested{One: false, Two: 2}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x01, 0x00}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify struct with includeIf handles relative reference and includes if true", func(t *testing.T) {
		type Nested struct {
			One bool
			Two uint8 `bcincludeif:"One"`
		}

		type StructUnderTest struct {
			One    bool
			Nested Nested
		}

		instance := &StructUnderTest{One: false, Nested: Nested{One: true, Two: 2}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x00, 0x01, 0x02}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify nested struct with includeIf handles absolute reference and includes if true", func(t *testing.T) {
		type Nested struct {
			One bool
			Two uint8 `bcincludeif:".Nested.One"`
		}

		type StructUnderTest struct {
			One    bool
			Nested Nested
		}

		instance := &StructUnderTest{One: false, Nested: Nested{One: true, Two: 2}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x00, 0x01, 0x02}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify struct with includeIf marshals a field which matches the integer condition", func(t *testing.T) {
		type StructUnderTest struct {
			One uint32
			Two uint8 `bcincludeif:".One==32"`
		}

		instance := &StructUnderTest{One: 32, Two: 2}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x20, 0x00, 0x00, 0x00, 0x02}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify struct with includeIf does not marshal a field which does not matche the integer condition", func(t *testing.T) {
		type StructUnderTest struct {
			One uint32
			Two uint8 `bcincludeif:".One==32"`
		}

		instance := &StructUnderTest{One: 16, Two: 2}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x10, 0x00, 0x00, 0x00}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("two 3 bit uints are written", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8 `bcfieldwidth:"3"`
			Two uint8 `bcfieldwidth:"3"`
		}

		instance := &StructUnderTest{One: 0b101, Two: 0b101}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0b10110100}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("writing a multibyte integer which is not a standard size", func(t *testing.T) {
		type StructUnderTest struct {
			One uint32 `bcfieldwidth:"24"`
		}

		instance := &StructUnderTest{One: 0x00aabbcc}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0xcc, 0xbb, 0xaa}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("writing an uint with a larger value than the field width results in error", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8 `bcfieldwidth:"3"`
		}

		instance := &StructUnderTest{One: 8}
		_, err := Marshal(instance)

		assert.Error(t, err)
	})

	t.Run("an error is thrown when a non byte aligned width over 8 bits is requested", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8 `bcfieldwidth:"9"`
		}

		instance := &StructUnderTest{One: 0}
		_, err := Marshal(instance)

		assert.Error(t, err)
	})

	t.Run("two 1 bit bools and a 6 bit uint are written", func(t *testing.T) {
		type StructUnderTest struct {
			One   bool  `bcfieldwidth:"1"`
			Two   bool  `bcfieldwidth:"1"`
			Three uint8 `bcfieldwidth:"6"`
		}

		instance := &StructUnderTest{One: true, Two: true, Three: 0x2d}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0b11101101}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("6 bit uint, array with 2 bit prefix is written", func(t *testing.T) {
		type StructUnderTest struct {
			One uint8  `bcfieldwidth:"6"`
			Two []byte `bcsliceprefix:"2"`
		}

		instance := &StructUnderTest{One: 0x2d, Two: []byte{0x00, 0x01}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0b10110110, 0x00, 0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("supports pointers which implement Marshaller interface", func(t *testing.T) {
		type StructUnderTest struct {
			One *CustomField
		}

		instance := &StructUnderTest{One: &CustomField{Value: 0}}
		actualBytes, err := Marshal(instance)

		expectedBytes := []byte{0x04, 'Z', 'E', 'R', 'O'}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})
}

type CustomField struct {
	Value uint8
}

func (f *CustomField) Marshal(bb *bitbuffer.BitBuffer) error {
	if f.Value == 0 {
		return MarshalToBitBuffer(bb, "ZERO")
	} else {
		return MarshalToBitBuffer(bb, "ONE")
	}
}

func (f *CustomField) Unmarshal(bb *bitbuffer.BitBuffer) error {
	s := ""
	pointer := &s

	err := UnmarshalFromBitBuffer(bb, pointer)

	if err != nil {
		return err
	}

	if s == "ZERO" {
		f.Value = 0
	} else {
		f.Value = 1
	}

	return nil
}
