package bytecodec

import (
	"testing"

	"github.com/shimmeringbee/bytecodec/bitbuffer"
	"github.com/stretchr/testify/assert"
)

func TestTagsEndian(t *testing.T) {
	t.Run("verifies that unannotated tag returns little endian", func(t *testing.T) {
		expectedValue := bitbuffer.LittleEndian
		actualValue := tagEndianness("")

		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated tag with little endian returns little endian", func(t *testing.T) {
		expectedValue := bitbuffer.LittleEndian
		actualValue := tagEndianness(`bcendian:"little"`)

		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated tag with big endian returns big endian", func(t *testing.T) {
		expectedValue := bitbuffer.BigEndian
		actualValue := tagEndianness(`bcendian:"big"`)
		assert.Equal(t, expectedValue, actualValue)
	})
}

func TestTagsArrayPrefix(t *testing.T) {
	t.Run("verifies that unannotated results in a length tag of 0 and little endian", func(t *testing.T) {
		expectedValue := SlicePrefixTag{
			Size:   0,
			Endian: bitbuffer.LittleEndian,
		}
		actualValue, err := tagSlicePrefix("")

		assert.NoError(t, err)
		assert.False(t, actualValue.HasPrefix())
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with length of one and no endian", func(t *testing.T) {
		expectedValue := SlicePrefixTag{
			Size:   8,
			Endian: bitbuffer.LittleEndian,
		}
		actualValue, err := tagSlicePrefix(`bcsliceprefix:"8"`)

		assert.NoError(t, err)
		assert.True(t, actualValue.HasPrefix())
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with length of one and no endian", func(t *testing.T) {
		expectedValue := SlicePrefixTag{
			Size:   16,
			Endian: bitbuffer.BigEndian,
		}
		actualValue, err := tagSlicePrefix(`bcsliceprefix:"16,big"`)

		assert.NoError(t, err)
		assert.True(t, actualValue.HasPrefix())
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verify that parse of invalid bit count returns error", func(t *testing.T) {
		_, err := tagSlicePrefix(`bcsliceprefix:"SPOON,big"`)

		assert.Error(t, err)
	})
}

func TestTagsString(t *testing.T) {
	t.Run("verifies that unannotated results in a prefix with a length tag of 0 and little endian", func(t *testing.T) {
		expectedValue := StringTypeTag{
			Termination: Prefix,
			Size:        8,
			Endian:      bitbuffer.LittleEndian,
		}
		actualValue, err := tagStringType("")

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with a prefix and length of 16 bits and big endian", func(t *testing.T) {
		expectedValue := StringTypeTag{
			Termination: Prefix,
			Size:        16,
			Endian:      bitbuffer.BigEndian,
		}
		actualValue, err := tagStringType(`bcstringtype:"prefix,16,big"`)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with null", func(t *testing.T) {
		expectedValue := StringTypeTag{
			Termination: Null,
			Size:        0,
			Endian:      bitbuffer.LittleEndian,
		}
		actualValue, err := tagStringType(`bcstringtype:"null"`)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with null with a padding size", func(t *testing.T) {
		expectedValue := StringTypeTag{
			Termination: Null,
			Size:        8,
			Endian:      bitbuffer.LittleEndian,
		}
		actualValue, err := tagStringType(`bcstringtype:"null,8"`)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, actualValue)
	})

	t.Run("verifies that annotated with null, and an invalid padding size", func(t *testing.T) {
		_, err := tagStringType(`bcstringtype:"null,SPOON,big"`)

		assert.Error(t, err)
	})

	t.Run("verifies that annotated with prefix, and an invalid padding size", func(t *testing.T) {
		_, err := tagStringType(`bcstringtype:"prefix,SPOON,big"`)

		assert.Error(t, err)
	})
}

func TestTagsIncludeIf(t *testing.T) {
	t.Run("verifies that single path non relative defaults to true", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:".field"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  false,
			FieldPath: []string{"field"},
			Operation: Equal,
			Value:     "",
		}, tag)
		assert.NoError(t, err)
	})

	t.Run("verifies that single path relative defaults to true", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:"field"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  true,
			FieldPath: []string{"field"},
			Operation: Equal,
			Value:     "",
		}, tag)
		assert.NoError(t, err)
	})

	t.Run("verifies that multiple path non relative defaults to true", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:".fieldOne.fieldTwo"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  false,
			FieldPath: []string{"fieldOne", "fieldTwo"},
			Operation: Equal,
			Value:     "",
		}, tag)
		assert.NoError(t, err)
	})

	t.Run("verifies that multiple path relative defaults to true", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:"fieldOne.fieldTwo"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  true,
			FieldPath: []string{"fieldOne", "fieldTwo"},
			Operation: Equal,
			Value:     "",
		}, tag)
		assert.NoError(t, err)
	})

	t.Run("verifies that single path non relative set to true", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:".field==true"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  false,
			FieldPath: []string{"field"},
			Operation: Equal,
			Value:     "true",
		}, tag)
		assert.NoError(t, err)
	})

	t.Run("verifies that single path non relative set to false", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:".field==false"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  false,
			FieldPath: []string{"field"},
			Operation: Equal,
			Value:     "false",
		}, tag)
		assert.NoError(t, err)
	})

	t.Run("verifies that single path non relative set to false", func(t *testing.T) {
		tag, err := tagIncludeIf(`bcincludeif:".field!=false"`)

		assert.Equal(t, IncludeIfTag{
			Relative:  false,
			FieldPath: []string{"field"},
			Operation: NotEqual,
			Value:     "false",
		}, tag)
		assert.NoError(t, err)
	})
}

func TestFieldWidthTag(t *testing.T) {
	t.Run("missing tag passes default width through", func(t *testing.T) {
		actualValue, err := tagFieldWidth(``)

		assert.NoError(t, err)
		assert.True(t, actualValue.Default)

		expectedValue := 16

		assert.Equal(t, expectedValue, actualValue.Width(expectedValue))
	})

	t.Run("provided tag overrides default width", func(t *testing.T) {
		actualValue, err := tagFieldWidth(`bcfieldwidth:"15"`)

		assert.NoError(t, err)
		assert.False(t, actualValue.Default)

		expectedValue := 15

		assert.Equal(t, expectedValue, actualValue.Width(8))
	})

	t.Run("tag with invalid bit count errors", func(t *testing.T) {
		_, err := tagFieldWidth(`bcfieldwidth:"SPOON"`)

		assert.Error(t, err)
	})
}
