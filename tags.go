package bytecodec

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/shimmeringbee/bytecodec/bitbuffer"
)

type StringTermination uint8

const (
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

func tagEndianness(tag reflect.StructTag) bitbuffer.Endian {
	if tag.Get(TagEndian) == BigEndianKeyword {
		return bitbuffer.BigEndian
	}

	return bitbuffer.LittleEndian
}

type SlicePrefixTag struct {
	Size   uint8
	Endian bitbuffer.Endian
}

func (l SlicePrefixTag) HasPrefix() bool {
	return l.Size > 0
}

func tagSlicePrefix(tag reflect.StructTag) (l SlicePrefixTag, err error) {
	l.Endian = bitbuffer.LittleEndian

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
		l.Endian = bitbuffer.BigEndian
	}

	return
}

type StringTypeTag struct {
	Termination StringTermination
	Size        uint8
	Endian      bitbuffer.Endian
}

func tagStringType(tag reflect.StructTag) (s StringTypeTag, err error) {
	s.Termination = Prefix
	s.Size = 8
	s.Endian = bitbuffer.LittleEndian

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
		s.Endian = bitbuffer.BigEndian
	}

	return
}

type IncludeIfOperation uint8

const (
	Equal    IncludeIfOperation = 0x00
	NotEqual IncludeIfOperation = 0x01
)

type IncludeIfTag struct {
	Relative  bool
	FieldPath []string

	Operation IncludeIfOperation

	Value string
}

var IncludeIfRegex = regexp.MustCompile(`^([a-zA-Z0-9.]+)(!=|==)?(.*)$`)

func tagIncludeIf(tag reflect.StructTag) (i IncludeIfTag, err error) {
	rawTag, tagPresent := tag.Lookup(TagIncludeIf)

	if !tagPresent {
		return IncludeIfTag{}, nil
	}

	matches := IncludeIfRegex.FindAllSubmatch([]byte(rawTag), -1)

	path := string(matches[0][1])
	operator := string(matches[0][2])

	i.Value = string(matches[0][3])
	i.Relative = path[0] != '.'

	pathParts := strings.Split(path, ".")
	partStart := 1

	if i.Relative {
		partStart = 0
	}

	i.FieldPath = pathParts[partStart:]

	if operator == "==" || operator == "" {
		i.Operation = Equal
	} else if operator == "!=" {
		i.Operation = NotEqual
	} else {
		return IncludeIfTag{}, fmt.Errorf("'%s' is not a valid includeIf operator", operator)
	}

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
