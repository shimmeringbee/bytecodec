package bytecodec

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshal(t *testing.T) {
	t.Run("verifies that error is thrown on unsupported type", func(t *testing.T) {
		type StructUnderTest struct {
			One chan bool
		}

		instance := &StructUnderTest{}
		_, err := Marshall(instance)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, UnsupportedType))
		assert.Equal(t, "unsupported type: field 'One' of type 'chan'", err.Error())
	})

	t.Run("verify byte and uint8 marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One byte
			Two uint8
		}

		instance := &StructUnderTest{One: 0x55, Two: 0xaa}
		actualBytes, err := Marshall(instance)

		expectedBytes := []byte{0x55, 0xaa}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify uint16 LE is marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One uint16 `bcendian:"little"`
		}

		instance := &StructUnderTest{One: 0x8001}
		actualBytes, err := Marshall(instance)

		expectedBytes := []byte{0x01, 0x80}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify uint16 BE is marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One uint16 `bcendian:"big"`
		}

		instance := &StructUnderTest{One: 0x8001}
		actualBytes, err := Marshall(instance)

		expectedBytes := []byte{0x80, 0x01}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify uint32 LE is marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One uint32 `bcendian:"little"`
		}

		instance := &StructUnderTest{One: 0x80010203}
		actualBytes, err := Marshall(instance)

		expectedBytes := []byte{0x03, 0x02, 0x01, 0x80}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify uint32 BE is marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One uint32 `bcendian:"big"`
		}

		instance := &StructUnderTest{One: 0x80010203}
		actualBytes, err := Marshall(instance)

		expectedBytes := []byte{0x80, 0x01, 0x02, 0x03}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify uint64 LE is marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One uint64 `bcendian:"little"`
		}

		instance := &StructUnderTest{One: 0x8001020304050607}
		actualBytes, err := Marshall(instance)

		expectedBytes := []byte{0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x80}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("verify uint64 BE is marshals", func(t *testing.T) {
		type StructUnderTest struct {
			One uint64 `bcendian:"big"`
		}

		instance := &StructUnderTest{One: 0x8001020304050607}
		actualBytes, err := Marshall(instance)

		expectedBytes := []byte{0x80, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}

		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, actualBytes)
	})
}
