package bitbuffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_UintRead(t *testing.T) {
	t.Run("read of 4 bits returns value", func(t *testing.T) {
		data := []byte{0b01011111}
		expectedValue := uint64(5)

		bb := NewBitBufferFromBytes(data)

		actualValue, err := bb.ReadUint(LittleEndian, 4)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("read of 9 bits returns error", func(t *testing.T) {
		data := []byte{0xff, 0xff}

		bb := NewBitBufferFromBytes(data)

		_, err := bb.ReadUint(LittleEndian, 9)
		assert.Error(t, err)
	})

	t.Run("read of 16 bits, little endian returns value", func(t *testing.T) {
		data := []byte{0xaa, 0xdd}
		expectedValue := uint64(0xddaa)

		bb := NewBitBufferFromBytes(data)

		actualValue, err := bb.ReadUint(LittleEndian, 16)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("read of 16 bits, big endian returns value", func(t *testing.T) {
		data := []byte{0xaa, 0xdd}
		expectedValue := uint64(0xaadd)

		bb := NewBitBufferFromBytes(data)

		actualValue, err := bb.ReadUint(BigEndian, 16)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})
}

func Test_UintWrite(t *testing.T) {
	t.Run("write of 4 bits returns value", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteUint(5, LittleEndian, 4)
		assert.NoError(t, err)

		expectedValue := []byte{0b01010000}
		actualValue := bb.Bytes()

		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("write of 9 bits returns error", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteUint(5, LittleEndian, 9)
		assert.Error(t, err)
	})

	t.Run("read of 16 bits, little endian returns value", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteUint(0xddaa, LittleEndian, 16)
		assert.NoError(t, err)

		expectedValue := []byte{0xaa, 0xdd}
		actualValue := bb.Bytes()

		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("read of 16 bits, big endian returns value", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteUint(0xddaa, BigEndian, 16)
		assert.NoError(t, err)

		expectedValue := []byte{0xdd, 0xaa}
		actualValue := bb.Bytes()

		assert.Equal(t, expectedValue, actualValue)
	})
}
