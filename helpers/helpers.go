package helpers

import "encoding/binary"

func GetByteFrom16(x uint16) byte {
	temp := make([]byte, 2)
	binary.LittleEndian.PutUint16(temp, x)
	return temp[0]
}

func GetByteFromInt(x int) byte {
	temp := make([]byte, 4)
	binary.LittleEndian.PutUint32(temp, uint32(x))
	return temp[0]
}
