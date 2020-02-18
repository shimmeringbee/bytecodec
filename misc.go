package bytecodec

import (
	"github.com/shimmeringbee/bytecodec/bitbuffer"
)

type Marshaler interface {
	Marshal(*bitbuffer.BitBuffer) error
}

type Unmarshaler interface {
	Unmarshal(*bitbuffer.BitBuffer) error
}
