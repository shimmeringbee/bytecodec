package bytecodec

type EndianTag uint8

type StringTermination uint8

const (
	BigEndian    EndianTag = 0
	LittleEndian EndianTag = 1

	Prefix StringTermination = 0
	Null   StringTermination = 1

	TagEndian = "bcendian"
	TagLength = "bclength"
	TagString = "bcstring"
)
