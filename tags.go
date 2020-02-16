package bytecodec

import (
	"reflect"
	"strconv"
	"strings"
)

type EndianTag uint8

type StringTermination uint8

const (
	BigEndian    EndianTag = 0
	LittleEndian EndianTag = 1

	Prefix StringTermination = 0
	Null   StringTermination = 1

	TagEndian      = "bcendian"
	TagSlicePrefix = "bcsliceprefix"
	TagStringType  = "bcstringtype"
	TagIncludeIf   = "bcincludeif"
	TagFieldWidth  = "bcfieldwidth"

	BigEndianKeyword       = "big"
	NullTerminationKeyword = "null"
)

func tagEndianness(tag reflect.StructTag) EndianTag {
	if tag.Get(TagEndian) == BigEndianKeyword {
		return BigEndian
	}

	return LittleEndian
}

type SlicePrefixTag struct {
	Size   uint8
	Endian EndianTag
}

func (l SlicePrefixTag) HasPrefix() bool {
	return l.Size > 0
}

func tagSlicePrefix(tag reflect.StructTag) (l SlicePrefixTag, err error) {
	l.Endian = LittleEndian

	rawTag, tagPresent := tag.Lookup(TagSlicePrefix)

	if !tagPresent {
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

	l.Size = uint8(bitCount)

	if len(splitTag) < 2 {
		return
	}

	if splitTag[1] == BigEndianKeyword {
		l.Endian = BigEndian
	}

	return
}

type StringTypeTag struct {
	Termination StringTermination
	Size        uint8
	Endian      EndianTag
}

func tagStringType(tag reflect.StructTag) (s StringTypeTag, err error) {
	s.Termination = Prefix
	s.Size = 8
	s.Endian = LittleEndian

	rawTag, tagPresent := tag.Lookup(TagStringType)

	if !tagPresent {
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
	rawTag, tagPresent := tag.Lookup(TagIncludeIf)

	if !tagPresent {
		return IncludeIfTag{}, nil
	}

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

func (i IncludeIfTag) HasIncludeIf() bool {
	return len(i.FieldPath) > 0
}

type FieldWidthTag struct {
	Default  bool
	BitWidth int
}

func (t FieldWidthTag) Width(defaultWidth int) int {
	if t.Default {
		return defaultWidth
	} else {
		return t.BitWidth
	}
}

func tagFieldWidth(tag reflect.StructTag) (t FieldWidthTag, err error) {
	rawTag, tagPresent := tag.Lookup(TagFieldWidth)

	if !tagPresent {
		t.Default = true
	} else {
		t.Default = false

		width, err := strconv.ParseInt(rawTag, 10, 8)

		if err != nil {
			return t, err
		}

		t.BitWidth = int(width)
	}

	return
}
