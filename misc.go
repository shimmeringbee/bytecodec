package bytecodec

import (
	"errors"
	"reflect"

	"github.com/shimmeringbee/bytecodec/bitbuffer"
)

var UnsupportedType = errors.New("unsupported type")

type Context struct {
	Root         reflect.Value
	CurrentIndex int
}

type Marshaler interface {
	Marshal(*bitbuffer.BitBuffer, Context) error
}

type Unmarshaler interface {
	Unmarshal(*bitbuffer.BitBuffer, Context) error
}
