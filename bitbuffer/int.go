package bitbuffer

import (
	"errors"
)

var ErrorNotByteMultiple = errors.New("can not read a non byte multiple int")

func checkByteMultiple(length int) error {
	if length%8 != 0 {
		return ErrorNotByteMultiple
	}

	return nil
}

func (bb *BitBuffer) WriteInt(value int64, endian Endian, length int) error {
	if err := checkByteMultiple(length); err != nil {
		return err
	}

	bytes := length / 8

	if endian == BigEndian {
		for i := 0; i < bytes; i++ {
			shiftOffset := (bytes - i - 1) * 8
			if err := bb.WriteByte(byte(value >> shiftOffset)); err != nil {
				return err
			}
		}
	} else {
		for i := 0; i < bytes; i++ {
			shiftOffset := i * 8
			if err := bb.WriteByte(byte(value >> shiftOffset)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (bb *BitBuffer) ReadInt(endian Endian, length int) (int64, error) {
	if err := checkByteMultiple(length); err != nil {
		return 0, err
	}

	readValue := int64(0)
	bytes := length / 8

	if endian == BigEndian {
		for i := 0; i < bytes; i++ {
			if b, err := bb.ReadByte(); err != nil {
				return 0, nil
			} else {
				shiftOffset := (bytes - i - 1) * 8
				readValue |= int64(int8(b)) << shiftOffset
			}
		}
	} else {
		for i := 0; i < bytes; i++ {
			if b, err := bb.ReadByte(); err != nil {
				return 0, nil
			} else {
				shiftOffset := i * 8
				readValue |= int64(int8(b)) << shiftOffset
			}
		}
	}

	return readValue, nil
}
