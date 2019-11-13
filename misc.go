package bytecodec

type Endian uint8

const (
	BigEndian    Endian = 0
	LittleEndian Endian = 1

	TagEndian = "bcendian"
	TagLength = "bclength"
)
