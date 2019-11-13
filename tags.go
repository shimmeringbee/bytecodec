package bytecodec

import (
	"reflect"
	"strconv"
	"strings"
)

const BigEndianKeyword string = "big"

func tagEndianness(tag reflect.StructTag) Endian {
	if tag.Get(TagEndian) == BigEndianKeyword {
		return BigEndian
	}

	return LittleEndian
}

type Length struct {
	Size   uint8
	Endian Endian
}

func (l Length) ShouldInsert() bool {
	return l.Size > 0
}

func tagLength(tag reflect.StructTag) (l Length, err error) {
	l.Endian = LittleEndian

	rawTag := tag.Get(TagLength)

	if rawTag == "" {
		return
	}

	splitTag := strings.Split(rawTag, ",")

	if len(splitTag) < 1 {
		return
	}

	bitCount, err := strconv.Atoi(splitTag[0])
	if err != nil {
		return
	}

	l.Size = uint8((bitCount + 7) / 8)

	if len(splitTag) < 2 {
		return
	}

	if splitTag[1] == BigEndianKeyword {
		l.Endian = BigEndian
	}

	return
}
