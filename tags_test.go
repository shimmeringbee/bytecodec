package bytecodec

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagsEndian(t *testing.T) {
	t.Run("verifies that unannoated tag returns little endian", func(t *testing.T) {
		expectedValue := LittleEndian
		actualValue := tagEndianness("")

		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated tag with little endian returns little endian", func(t *testing.T) {
		expectedValue := LittleEndian
		actualValue := tagEndianness(`bcendian:"little"`)

		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated tag with big endian returns big endian", func(t *testing.T) {
		expectedValue := BigEndian
		actualValue := tagEndianness(`bcendian:"big"`)

		assert.Equal(t, expectedValue, actualValue)
	})
}
