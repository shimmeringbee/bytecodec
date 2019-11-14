package bytecodec

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshall(t *testing.T) {
	t.Run("verifies that error is thrown on unsupported type", func(t *testing.T) {
		type StructUnderTest struct {
			One chan bool
		}

		var data []byte
		instance := StructUnderTest{}
		err := Unmarshall(data, instance)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, UnsupportedType))
		assert.Equal(t, "unsupported type: field 'One' of type 'chan'", err.Error())
	})

	t.Run("verify byte and uint8 unmarshals", func(t *testing.T) {
		type StructUnderTest struct {
			One byte
			Two uint8
		}

		expectedStruct := StructUnderTest{One: 0x55, Two: 0xaa}
		data, _ := Marshall(expectedStruct)

		actualStruct := StructUnderTest{}
		err := Unmarshall(data, actualStruct)

		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, actualStruct)
	})
}
