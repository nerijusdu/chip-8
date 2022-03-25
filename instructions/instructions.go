package instructions

import (
	. "chip8/models"
	"fmt"
)

func ClearScreen(data *GameData) {
	for i := 0; i < data.PixelCount; i++ {
		data.Pixels[i] = false
	}
}

func Jump(data *GameData, nnn uint16) {
	data.PC = nnn
}

func SetRegister(data *GameData, x uint16, nn byte) {
	data.Registers[x] = nn
}

func AddValueToRegister(data *GameData, x uint16, nn byte) {
	data.Registers[x] += nn
}

func SetIndexRegister(data *GameData, nnn uint16) {
	data.I = nnn
}

func DisplayDraw(data *GameData, x uint16, y uint16, n uint16) {
	X := data.Registers[x] & 63
	Y := data.Registers[y] & 31
	data.Registers[0xF] = 0

	var i uint16 = 0
	for ; i < n; i++ {
		sprite := data.Memory[data.I+i]
		bits := fmt.Sprintf("%08b", sprite)
		yy := uint16(Y) + uint16(i)
		if yy > 31 {
			break
		}

		for j, bit := range bits {
			isOn := bit == '1'
			xx := uint16(X) + uint16(j)
			if xx > 63 {
				break
			}

			currentPixel := data.Pixels.Get(xx, yy)
			if isOn && currentPixel {
				data.Registers[0xF] = 1
			}
			data.Pixels.Set(xx, yy, isOn && !currentPixel)
		}
	}
}
