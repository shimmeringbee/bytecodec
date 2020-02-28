package bitbuffer

import (
	"errors"
	"math"
	"strings"
)

var ErrorStringTooLarge = errors.New("string length exceeds capacity of length prefix")

func (bb *BitBuffer) WriteStringNullTerminated(data string, paddedLength int) error {
	dataLen := len(data)

	if paddedLength > 0 && dataLen >= paddedLength {
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
	maxLength := int(math.Pow(2, float64(length)) - 1)

	if len(data) > maxLength {
		return ErrorStringTooLarge
	}

	if err := bb.WriteUint(uint64(len(data)), endian, length); err != nil {
		return err
	}

	return bb.WriteString(data)
}

func (bb *BitBuffer) ReadStringNullTerminated(paddedLength int) (string, error) {
	sb := strings.Builder{}

	maxLength := math.MaxInt64

	if paddedLength != 0 {
		maxLength = paddedLength
	}

	i := 0

	readByte := ^byte(0)

	for ; i < maxLength; i++ {
		rB, err := bb.ReadByte()
		if err != nil {
			return "", err
		}

		readByte = rB

		if readByte == 0 {
			i++
			break
		}

		sb.WriteByte(readByte)
	}

	if readByte != 0 {
		return "", errors.New("no null termination found in string")
	}

	for ; i < paddedLength; i++ {
		_, err := bb.ReadByte()
		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

func (bb *BitBuffer) ReadStringLengthPrefixed(endian Endian, length int) (string, error) {
	sb := strings.Builder{}

	stringLength, err := bb.ReadUint(endian, length)
	if err != nil {
		return "", err
	}

	for i := 0; i < int(stringLength); i++ {
		readByte, err := bb.ReadByte()
		if err != nil {
			return "", err
		}

		sb.WriteByte(readByte)
	}

	return sb.String(), nil
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
