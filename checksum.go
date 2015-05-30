package typeframed

import (
	"hash/crc32"
)

type Checksum func(bytes []byte) []byte

func Crc32Checksum(bytes []byte) []byte {
	chksum := crc32.ChecksumIEEE(bytes)
	buf := make([]byte, 8)
	buf[0] = byte(chksum >> 56)
	buf[1] = byte(chksum >> 48)
	buf[2] = byte(chksum >> 40)
	buf[3] = byte(chksum >> 32)
	buf[4] = byte(chksum >> 24)
	buf[5] = byte(chksum >> 16)
	buf[6] = byte(chksum >> 8)
	buf[7] = byte(chksum >> 0)
	return buf
}

type CorruptedChecksumError struct {
    Msg string
}

func (e *CorruptedChecksumError) Error() string {
    return e.Msg
}

func NewCorruptedChecksum() error {
	return &CorruptedChecksumError{"Corrupted checksum"}
} 