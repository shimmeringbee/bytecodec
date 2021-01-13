package bitbuffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BitsBytes(t *testing.T) {
	t.Run("attempting to write more than 8 bits errors", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteBits(0x00, 9)
		assert.Error(t, err)
	})

	t.Run("writing bytes works normally", func(t *testing.T) {
		bb := NewBitBuffer()

		bb.WriteByte(0xaa)
		bb.WriteByte(0xdd)

		actualBytes := bb.Bytes()
		expectedBytes := []byte{0xaa, 0xdd}

		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("writing two 2 bit and then one 8 bit", func(t *testing.T) {
		bb := NewBitBuffer()

		bb.WriteBits(0x01, 2)
		bb.WriteBits(0x02, 2)
		bb.WriteByte(0xaa)

		actualBytes := bb.Bytes()
		expectedBytes := []byte{0x6a, 0xa0}

		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("attempting to read more than 8 bits errors", func(t *testing.T) {
		bb := NewBitBuffer()

		_, err := bb.ReadBits(9)
		assert.Error(t, err)
	})

	t.Run("reading bytes works normally", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0xaa, 0xdd})

		one, _ := bb.ReadByte()
		two, _ := bb.ReadByte()
		assert.Equal(t, byte(0xaa), one)
		assert.Equal(t, byte(0xdd), two)
	})

	t.Run("reading two 2 bit and then one 8 bit", func(t *testing.T) {
		bb := NewBitBufferFromBytes([]byte{0x6a, 0xa0})

		one, _ := bb.ReadBits(2)
		two, _ := bb.ReadBits(2)
		three, _ := bb.ReadByte()

		assert.Equal(t, byte(0x01), one)
		assert.Equal(t, byte(0x02), two)
		assert.Equal(t, byte(0xaa), three)
	})
}
