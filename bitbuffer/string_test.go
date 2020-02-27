package bitbuffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

	t.Run("write null terminated string, padding, full length string", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteStringNullTerminated("HiHi", 4)
		assert.NoError(t, err)

		actualBytes := bb.Bytes()
		expectedBytes := []byte{0x48, 0x69, 0x48, 0x69}

		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("write null terminated string, padding, string longer than padding", func(t *testing.T) {
		bb := NewBitBuffer()

		err := bb.WriteStringNullTerminated("HiHi", 2)
		assert.Error(t, err)
	})

	t.Run("writing a string", func(t *testing.T) {
		bb := NewBitBuffer()

		bb.WriteString("Hi")

		actualBytes := bb.Bytes()
		expectedBytes := []byte{0x48, 0x69}

		assert.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("writing a unicode string", func(t *testing.T) {
		bb := NewBitBuffer()

		bb.WriteString("🤬")

		actualBytes := bb.Bytes()
		expectedBytes := []byte{0xf0, 0x9f, 0xa4, 0xac}

		assert.Equal(t, expectedBytes, actualBytes)
	})
}
