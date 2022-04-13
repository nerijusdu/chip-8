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

func GetKey(key string) byte {
	switch key {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "0":
		return 0
	case "A":
		return 0xA
	case "B":
		return 0xB
	case "C":
		return 0xC
	case "D":
		return 0xD
	case "E":
		return 0xE
	case "F":
		return 0xF
	}
	return 0
}
