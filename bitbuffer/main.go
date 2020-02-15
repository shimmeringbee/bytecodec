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

func (bb *BitBuffer) ReadByte() (byte, error) {
	return bb.ReadBits(8)
}

func (bb *BitBuffer) WriteByte(c byte) error {
	return bb.WriteBits(c, 8)
}

func (bb *BitBuffer) ReadBits(bitCount int) (byte, error) {
	return 0, nil
}

func (bb *BitBuffer) WriteBits(bits byte, bitCount int) error {
	if bitCount > maxBitOperations {
		return ErrorTooManyBitsInOperation
	}

	mask := byte(1 << (bitCount - 1))

	for i := 0; i < bitCount; i++ {
		bit := (bits & mask) == mask
		bits = bits << 1

		bb.unhandled = bb.unhandled << 1

		if bit {
			bb.unhandled = bb.unhandled | 0x01
		}

		bb.offset++

		if bb.offset == 8 {
			_ = bb.buf.WriteByte(bb.unhandled)
			bb.unhandled = 0
			bb.offset = 0
		}
	}

	return nil
}
