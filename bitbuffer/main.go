package bitbuffer

import (
	"bytes"
	"errors"
)

const maxBitOperations int = 8

var ErrorTooManyBitsInOperation = errors.New("bit buffer can only perform operations on 8 or fewer bits")

func NewBitBuffer() *BitBuffer {
	return &BitBuffer{
		buf: &bytes.Buffer{},
	}
}

func NewBitBufferFromBytes(data []byte) *BitBuffer {
	return &BitBuffer{
		buf: bytes.NewBuffer(data),
	}
}

type BitBuffer struct {
	buf       *bytes.Buffer
	unhandled byte
	offset    uint8
}

func (bb *BitBuffer) Bytes() []byte {
	if bb.offset != 0 {
		_ = bb.WriteBits(0, int(8-bb.offset))
	}
	return bb.buf.Bytes()
}
