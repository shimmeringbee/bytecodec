package bitbuffer

import (
	"fmt"
	"math"
)

func (bb *BitBuffer) ReadUint(endian Endian, length int) (uint64, error) {
	if length < 8 {
		data, err := bb.ReadBits(length)
		return uint64(data), err
	}

	if length%8 != 0 {
		return 0, fmt.Errorf("unable to handle arbitary bit widths > 8 bits, %d requested", length)
	}

	readValue := uint64(0)
	bytes := length / 8

	switch endian {
	case BigEndian:
		for i := 0; i < bytes; i++ {
			readByte, err := bb.ReadByte()
			if err != nil {
				return 0, err
			}

			shiftOffset := (bytes - i - 1) * 8
			readValue |= uint64(readByte) << shiftOffset
		}
	case LittleEndian:
		for i := 0; i < bytes; i++ {
			readByte, err := bb.ReadByte()
			if err != nil {
				return 0, err
			}

			shiftOffset := i * 8
			readValue |= uint64(readByte) << shiftOffset
		}
	}

	return readValue, nil
}

func (bb *BitBuffer) WriteUint(value uint64, endian Endian, length int) error {
	maxValue := math.Pow(2, float64(length)) - 1

	if float64(value) > maxValue {
		return fmt.Errorf("cannot marshal value %d into %d bit field", value, length)
	}

	if length < 8 {
		if err := bb.WriteBits(byte(value), length); err != nil {
			return err
		}
	} else {
		bytes := length / 8

		if length%8 != 0 {
			return fmt.Errorf("unable to handle arbitary bit widths > 8 bits, %d requested", length)
		}

		switch endian {
		case BigEndian:
			for i := 0; i < bytes; i++ {
				shiftOffset := (bytes - i - 1) * 8
				if err := bb.WriteByte(byte(value >> shiftOffset)); err != nil {
					return err
				}
			}
		case LittleEndian:
			for i := 0; i < bytes; i++ {
				shiftOffset := i * 8
				if err := bb.WriteByte(byte(value >> shiftOffset)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
