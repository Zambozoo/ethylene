package core

import "encoding/binary"

func LongBytes(i uint64) []byte {
	return binary.LittleEndian.AppendUint64(nil, i)
}

func IntBytes(i int) []byte {
	return binary.LittleEndian.AppendUint32(nil, uint32(i))
}
