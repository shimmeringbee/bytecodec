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
		expectedValue := LengthTag{
			Size:   0,
			Endian: LittleEndian,
		}
		actualValue, err := tagLength("")

		assert.NoError(t, err)
		assert.False(t, actualValue.HasLength())
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with length of one and no endian", func(t *testing.T) {
		expectedValue := LengthTag{
			Size:   1,
			Endian: LittleEndian,
		}
		actualValue, err := tagLength(`bclength:"8"`)

		assert.NoError(t, err)
		assert.True(t, actualValue.HasLength())
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with length of one and no endian", func(t *testing.T) {
		expectedValue := LengthTag{
			Size:   2,
			Endian: BigEndian,
		}
		actualValue, err := tagLength(`bclength:"16,big"`)

		assert.NoError(t, err)
		assert.True(t, actualValue.HasLength())
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verify that parse of invalid bit count returns error", func(t *testing.T) {
		_, err := tagLength(`bclength:"SPOON,big"`)

		assert.Error(t, err)
	})
}

func TestTagsString(t *testing.T) {
	t.Run("verifies that unannotated results in a prefix with a length tag of 0 and little endian", func(t *testing.T) {
		expectedValue := StringTag{
			Termination: Prefix,
			Size:        8,
			Endian:      LittleEndian,
		}
		actualValue, err := tagString("")

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with a prefix and length of 16 bits and big endian", func(t *testing.T) {
		expectedValue := StringTag{
			Termination: Prefix,
			Size:        16,
			Endian:      BigEndian,
		}
		actualValue, err := tagString(`bcstring:"prefix,16,big"`)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with null", func(t *testing.T) {
		expectedValue := StringTag{
			Termination: Null,
			Size:        0,
			Endian:      LittleEndian,
		}
		actualValue, err := tagString(`bcstring:"null"`)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with null with a padding size", func(t *testing.T) {
		expectedValue := StringTag{
			Termination: Null,
			Size:        8,
			Endian:      LittleEndian,
		}
		actualValue, err := tagString(`bcstring:"null,8"`)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with null, and an invalid padding size", func(t *testing.T) {
		_, err := tagString(`bcstring:"null,SPOON,big"`)

		assert.Error(t, err)
	})

	t.Run("verifies that annotated with prefix, and an invalid padding size", func(t *testing.T) {
		_, err := tagString(`bcstring:"prefix,SPOON,big"`)

		assert.Error(t, err)
	})
}

func TestTagsIncludeIf(t *testing.T) {
	t.Run("verifies that single path non relative defaults to true", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:".field"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  false,
			FieldPath: []string{"field"},
			Value:     true,
		}, tag)
		assert.NoError(t, err)
	})

	t.Run("verifies that single path relative defaults to true", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:"field"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  true,
			FieldPath: []string{"field"},
			Value:     true,
		}, tag)
		assert.NoError(t, err)
	})

	t.Run("verifies that multiple path non relative defaults to true", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:".fieldOne.fieldTwo"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  false,
			FieldPath: []string{"fieldOne", "fieldTwo"},
			Value:     true,
		}, tag)
		assert.NoError(t, err)
	})

	t.Run("verifies that multiple path relative defaults to true", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:"fieldOne.fieldTwo"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  true,
			FieldPath: []string{"fieldOne", "fieldTwo"},
			Value:     true,
		}, tag)
		assert.NoError(t, err)
	})

	t.Run("verifies that single path non relative set to true", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:".field=true"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  false,
			FieldPath: []string{"field"},
			Value:     true,
		}, tag)
		assert.NoError(t, err)
	})

	t.Run("verifies that single path non relative set to false", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:".field=false"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  false,
			FieldPath: []string{"field"},
			Value:     false,
		}, tag)
		assert.NoError(t, err)
	})
}
