package bitbuffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	t.Run("write null terminated string, no padding", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteStringNullTerminated("Hi", 0)
		assert.NoError(t, err)

		actualBytes := bb.Bytes()
		expectedBytes := []byte{0x48, 0x69, 0x00}

		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("write null terminated string, padding, short string", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteStringNullTerminated("Hi", 4)
		assert.NoError(t, err)

		actualBytes := bb.Bytes()
		expectedBytes := []byte{0x48, 0x69, 0x00, 0x00}

		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("write null terminated string, padding, full length string (with null)", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteStringNullTerminated("HiH", 4)
		assert.NoError(t, err)

		actualBytes := bb.Bytes()
		expectedBytes := []byte{0x48, 0x69, 0x48, 0x00}

		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("write null terminated string, padding, string longer than padding", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteStringNullTerminated("HiHi", 2)
		assert.Error(t, err)
	})

	t.Run("write length prefixed string, 4 bits", func(t *testing.T) {
		bb := NewBitBuffer()
		bb.WriteBits(0b1111, 4)

		expectedBytes := []byte{0b11110010, 'H', 'i'}

		err := bb.WriteStringLengthPrefixed("Hi", LittleEndian, 4)
		assert.NoError(t, err)

		actualBytes := bb.Bytes()

		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("write length prefixed string, 2 bits, string too long", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteStringLengthPrefixed("Hello", LittleEndian, 2)
		assert.Error(t, err)
	})

	t.Run("write length prefixed string, 16 bits, little endian", func(t *testing.T) {
		bb := NewBitBuffer()
		expectedBytes := []byte{0x02, 0x00, 'H', 'i'}

		err := bb.WriteStringLengthPrefixed("Hi", LittleEndian, 16)
		assert.NoError(t, err)

		actualBytes := bb.Bytes()

		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("write length prefixed string, 16 bits, big endian", func(t *testing.T) {
		bb := NewBitBuffer()
		expectedBytes := []byte{0x00, 0x02, 'H', 'i'}

		err := bb.WriteStringLengthPrefixed("Hi", BigEndian, 16)
		assert.NoError(t, err)

		actualBytes := bb.Bytes()

		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("unmarshal null terminated string, no padding", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{'H', 'i', 0x00, 0x01})

		expectedString := "Hi"
		actualString, err := bb.ReadStringNullTerminated(0)

		assert.NoError(t, err)
		assert.Equal(t, expectedString, actualString)

		endCheck, err := bb.ReadUint(LittleEndian, 8)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0x01), endCheck)
	})

	t.Run("unmarshal null terminated string, 4 padding", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{'H', 'i', 0x00, 0x00, 0x01})

		expectedString := "Hi"
		actualString, err := bb.ReadStringNullTerminated(4)

		assert.NoError(t, err)
		assert.Equal(t, expectedString, actualString)

		endCheck, err := bb.ReadUint(LittleEndian, 8)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0x01), endCheck)
	})

	t.Run("unmarshal length prefixed string, 4 bits", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0b11110010, 'H', 'i'})

		expectedString := "Hi"
		bb.ReadBits(4)
		actualString, err := bb.ReadStringLengthPrefixed(LittleEndian, 4)

		assert.NoError(t, err)
		assert.Equal(t, expectedString, actualString)
	})

	t.Run("unmarshal length prefixed string, 16 bits, little endian", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0x02, 0x00, 'H', 'i'})

		expectedString := "Hi"
		actualString, err := bb.ReadStringLengthPrefixed(LittleEndian, 16)

		assert.NoError(t, err)
		assert.Equal(t, expectedString, actualString)
	})

	t.Run("unmarshal length prefixed string, 16 bits, big endian", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0x00, 0x02, 'H', 'i'})

		expectedString := "Hi"
		actualString, err := bb.ReadStringLengthPrefixed(BigEndian, 16)

		assert.NoError(t, err)
		assert.Equal(t, expectedString, actualString)
	})

	t.Run("writing a string", func(t *testing.T) {
		bb := NewBitBuffer()

		bb.writeString("Hi")

		actualBytes := bb.Bytes()
		expectedBytes := []byte{0x48, 0x69}

		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("writing a unicode string", func(t *testing.T) {
		bb := NewBitBuffer()

		bb.writeString("🤬")

		actualBytes := bb.Bytes()
		expectedBytes := []byte{0xf0, 0x9f, 0xa4, 0xac}

		assert.Equal(t, expectedBytes, actualBytes)
	})
}
