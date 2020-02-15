package bitbuffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BitBuffer(t *testing.T) {
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
}
