package bitbuffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadInt(t *testing.T) {
	t.Run("an error is thrown attempting to read none byte multiple int", func(t *testing.T) {
		bb := NewBitBuffer()

		_, err := bb.ReadInt(LittleEndian, 5)

		assert.Error(t, err)
	})

	t.Run("reading a 4 byte big endian positive integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0x7F, 0x00, 0x00, 0x00})

		expectedValue := int64(0x7f000000)
		actualValue, err := bb.ReadInt(BigEndian, 32)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("reading a 4 byte little endian positive integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0x00, 0x00, 0x00, 0x7f})

		expectedValue := int64(0x7f000000)
		actualValue, err := bb.ReadInt(LittleEndian, 32)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("reading a 4 byte big endian negative integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0xFF, 0xFF, 0xFF, 0x00})

		expectedValue := int64(-256)
		actualValue, err := bb.ReadInt(BigEndian, 32)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("reading a 4 byte little endian negative integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0x00, 0xFF, 0xFF, 0xFF})

		expectedValue := int64(-256)
		actualValue, err := bb.ReadInt(LittleEndian, 32)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("reading a 1 byte no endian positive integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0xFF})

		expectedValue := int64(-1)
		actualValue, err := bb.ReadInt(BigEndian, 8)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("reading a 1 byte no endian negative integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0x7F})

		expectedValue := int64(127)
		actualValue, err := bb.ReadInt(LittleEndian, 8)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("reading a 2 byte little endian negative integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0xb8, 0x09})

		expectedValue := int64(2488)
		actualValue, err := bb.ReadInt(LittleEndian, 16)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("reading a 2 byte big endian negative integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0x09, 0xb8})

		expectedValue := int64(2488)
		actualValue, err := bb.ReadInt(BigEndian, 16)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})
}

func Test_WriteInt(t *testing.T) {
	t.Run("an error is thrown attempting to write none byte multiple int", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteInt(0, LittleEndian, 5)

		assert.Error(t, err)
	})

	t.Run("writing a 4 byte big endian positive integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{})

		inputValue := int64(0x7f000000)
		expectedValue := []byte{0x7F, 0x00, 0x00, 0x00}

		err := bb.WriteInt(inputValue, BigEndian, 32)
		actualValue := bb.Bytes()

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("writing a 4 byte little endian positive integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{})

		inputValue := int64(0x7f000000)
		expectedValue := []byte{0x00, 0x00, 0x00, 0x7f}

		err := bb.WriteInt(inputValue, LittleEndian, 32)
		actualValue := bb.Bytes()

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("writing a 4 byte big endian negative integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{})

		inputValue := int64(-256)
		expectedValue := []byte{0xFF, 0xFF, 0xFF, 0x00}

		err := bb.WriteInt(inputValue, BigEndian, 32)
		actualValue := bb.Bytes()

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("writing a 4 byte little endian negative integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{})

		inputValue := int64(-256)
		expectedValue := []byte{0x00, 0xFF, 0xFF, 0xFF}

		err := bb.WriteInt(inputValue, LittleEndian, 32)
		actualValue := bb.Bytes()

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("writing a 1 byte no endian positive integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{})

		inputValue := int64(127)
		expectedValue := []byte{0x7F}

		err := bb.WriteInt(inputValue, BigEndian, 8)
		actualValue := bb.Bytes()

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("writing a 1 byte no endian negative integer", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{})

		inputValue := int64(-1)
		expectedValue := []byte{0xFF}

		err := bb.WriteInt(inputValue, LittleEndian, 8)
		actualValue := bb.Bytes()

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})
}
