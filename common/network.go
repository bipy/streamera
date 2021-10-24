package common

const (
	TimeStampPacketSize     = 8
	ContentLengthPacketSize = 8
	HeadPacketSize          = TimeStampPacketSize + ContentLengthPacketSize
)
