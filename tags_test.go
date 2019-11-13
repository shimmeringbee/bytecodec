package bytecodec

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagsEndian(t *testing.T) {
	t.Run("verifies that unannotated tag returns little endian", func(t *testing.T) {
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

func TestTagsLength(t *testing.T) {
	t.Run("verifies that unannotated results in a length tag of 0 and little endian", func(t *testing.T) {
		expectedValue := Length{
			Size:   0,
			Endian: LittleEndian,
		}
		actualValue, err := tagLength("")

		assert.NoError(t, err)
		assert.False(t, actualValue.ShouldInsert())
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with length of one and no endian", func(t *testing.T) {
		expectedValue := Length{
			Size:   1,
			Endian: LittleEndian,
		}
		actualValue, err := tagLength(`bclength:"8"`)

		assert.NoError(t, err)
		assert.True(t, actualValue.ShouldInsert())
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with length of one and no endian", func(t *testing.T) {
		expectedValue := Length{
			Size:   2,
			Endian: BigEndian,
		}
		actualValue, err := tagLength(`bclength:"16,big"`)

		assert.NoError(t, err)
		assert.True(t, actualValue.ShouldInsert())
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verify that parse of invalid bit count returns error", func(t *testing.T) {
		_, err := tagLength(`bclength:"SPOON,big"`)

		assert.Error(t, err)
	})
}
