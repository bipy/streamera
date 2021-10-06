package common

const (
    Width  = 1280
    Height = 720
    Depth  = 8
    Rgb    = 3
    Mode   = (Depth >> 3) * Rgb
)

const (
    TimeStampPacketSize = 8
    FramePacketSize  = Width * Height * Mode
    ChunkSize           = TimeStampPacketSize + FramePacketSize
)
