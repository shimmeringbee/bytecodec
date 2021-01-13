package bitbuffer

func (bb *BitBuffer) ReadByte() (byte, error) {
	return bb.ReadBits(8)
}

func (bb *BitBuffer) WriteByte(c byte) error {
	return bb.WriteBits(c, 8)
}

func (bb *BitBuffer) ReadBits(bitCount int) (byte, error) {
	if bitCount > maxBitOperations {
		return 0, ErrorTooManyBitsInOperation
	}

	if bb.offset == 0 && bitCount == 8 {
		return bb.buf.ReadByte()
	}

	retVal := byte(0)

	for i := 0; i < bitCount; i++ {
		if bb.offset == 0 {
			unhandled, err := bb.buf.ReadByte()

			if err != nil {
				return 0, err
			}

			bb.unhandled = unhandled
			bb.offset = 8
		}

		bit := bb.unhandled&0x80 == 0x80
		bb.unhandled <<= 1
		bb.offset--

		retVal <<= 1

		if bit {
			retVal |= 1
		}
	}

	return retVal, nil
}

func (bb *BitBuffer) WriteBits(bits byte, bitCount int) error {
	if bitCount > maxBitOperations {
		return ErrorTooManyBitsInOperation
	}

	if bb.offset == 0 && bitCount == 8 {
		return bb.buf.WriteByte(bits)
	}

	mask := byte(1 << (bitCount - 1))

	for i := 0; i < bitCount; i++ {
		bit := (bits & mask) == mask
		bits <<= 1

		bb.unhandled <<= 1

		if bit {
			bb.unhandled |= 0x01
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
