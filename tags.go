package bytecodec

import (
	"reflect"
	"strconv"
	"strings"
)

const BigEndianKeyword string = "big"
const NullTerminationKeyword string = "null"

func tagEndianness(tag reflect.StructTag) EndianTag {
	if tag.Get(TagEndian) == BigEndianKeyword {
		return BigEndian
	}

	return LittleEndian
}

type LengthTag struct {
	Size   uint8
	Endian EndianTag
}

func (l LengthTag) HasLength() bool {
	return l.Size > 0
}

func tagLength(tag reflect.StructTag) (l LengthTag, err error) {
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

type StringTag struct {
	Termination StringTermination
	Size        uint8
	Endian      EndianTag
}

func tagString(tag reflect.StructTag) (s StringTag, err error) {
	s.Termination = Prefix
	s.Size = 8
	s.Endian = LittleEndian

	rawTag := tag.Get(TagString)

	if rawTag == "" {
		return
	}

	splitTag := strings.Split(rawTag, ",")

	if splitTag[0] == NullTerminationKeyword {
		s.Termination = Null
		s.Size = 0
	}

	if len(splitTag) <= 1 {
		return
	}

	count, err := strconv.Atoi(splitTag[1])
	if err != nil {
		return
	}

	s.Size = uint8(count)

	if len(splitTag) <= 2 {
		return
	}

	if splitTag[2] == BigEndianKeyword {
		s.Endian = BigEndian
	}

	return
}

type IncludeIfTag struct {
	Relative  bool
	FieldPath []string
	Value     bool
}

func tagIncludeIf(tag reflect.StructTag) (i IncludeIfTag, err error) {
	rawTag := tag.Get(TagIncludeIf)

	i.Value = true
	i.Relative = rawTag[0] != '.'

	tagParts := strings.Split(rawTag, "=")

	if len(tagParts) > 1 {
		i.Value, err = strconv.ParseBool(tagParts[1])

		if err != nil {
			return
		}
	}

	pathParts := strings.Split(tagParts[0], ".")
	partStart := 1

	if i.Relative {
		partStart = 0
	}

	i.FieldPath = pathParts[partStart:]

	return
}
