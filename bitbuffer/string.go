package bitbuffer

import "errors"

func (bb *BitBuffer) WriteStringNullTerminated(data string, paddedLength int) error {
	dataLen := len(data)

	if paddedLength > 0 && dataLen > paddedLength {
		return errors.New("string to too long to write into padded string")
	}

	if err := bb.WriteString(data); err != nil {
		return err
	}

	if paddedLength > 0 {
		for i := 0; i < paddedLength-dataLen; i++ {
			if err := bb.WriteByte(0); err != nil {
				return err
			}
		}
	} else {
		return bb.WriteByte(0)
	}

	return nil
}

func (bb *BitBuffer) WriteStringLengthPrefixed(data string, endian Endian, length int) error {
	return nil
}

func (bb *BitBuffer) ReadStringNullTerminated(paddedLength int) (string, error) {
	return "", nil
}

func (bb *BitBuffer) ReadStringLengthPrefixed(endian Endian, length int) (string, error) {
	return "", nil
}

func (bb *BitBuffer) WriteString(data string) error {
	dataBytes := []byte(data)

	for _, b := range dataBytes {
		if err := bb.WriteByte(b); err != nil {
			return nil
		}
	}

	return nil
}
