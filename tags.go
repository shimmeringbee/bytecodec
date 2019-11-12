package bytecodec

import "reflect"

func tagEndianness(tag reflect.StructTag) Endian {
	if tag.Get(TagEndian) == "big" {
		return BigEndian
	}

	return LittleEndian
}
