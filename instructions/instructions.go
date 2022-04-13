package instructions

import (
	"chip8/helpers"
	. "chip8/models"
	"fmt"
	"math/rand"
)

func ClearScreen(data *GameData) {
	for i := 0; i < data.PixelCount; i++ {
		data.Pixels[i] = false
	}
}

func Jump(data *GameData, nnn uint16) {
	data.PC = nnn
}

func JumpWithOffset(data *GameData, nnn uint16) {
	data.PC = nnn + uint16(data.Registers[0])
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

func StartSubroutine(data *GameData, nnn uint16) {
	data.Stack.Push(data.PC)
	Jump(data, nnn)
}

func ReturnFromSubroutine(data *GameData) {
	pc := data.Stack.Pop()
	Jump(data, pc)
}

func Skip(data *GameData) {
	data.PC += 2
}

func SkipIfEqual(data *GameData, x uint16, nn byte, isEqualOptional ...bool) {
	isEqual := true
	if len(isEqualOptional) > 0 {
		isEqual = isEqualOptional[0]
	}

	result := data.Registers[x] == nn
	if result == isEqual {
		Skip(data)
	}
}

func SkipIfEqualRegisters(data *GameData, x uint16, y uint16, isEqualOptional ...bool) {
	isEqual := true
	if len(isEqualOptional) > 0 {
		isEqual = isEqualOptional[0]
	}

	result := data.Registers[x] == data.Registers[y]
	if result == isEqual {
		Skip(data)
	}
}

func AddWithCarry(data *GameData, x uint16, y uint16) {
	vx := uint16(data.Registers[x])
	vy := uint16(data.Registers[y])
	var vf byte = 0
	if vx+vy > 255 {
		vf = 1
	}

	data.Registers[0xF] = vf
	data.Registers[x] = data.Registers[x] + data.Registers[y]
}

func SubtractWithCarry(data *GameData, x uint16, a uint16, b uint16) {
	var vf byte = 0
	aa := data.Registers[a]
	bb := data.Registers[b]

	if a > b {
		vf = 1
	}

	data.Registers[x] = aa - bb
	data.Registers[0xF] = vf
}

func ShiftRight(data *GameData, x uint16) {
	var cf byte
	if (data.Registers[x] & 0x01) == 0x01 {
		cf = 1
	}
	data.Registers[0xF] = cf
	data.Registers[x] = data.Registers[x] / 2
}

func ShiftLeft(data *GameData, x uint16) {
	var cf byte
	if (data.Registers[x] & 0x80) == 0x80 {
		cf = 1
	}
	data.Registers[0xF] = cf
	data.Registers[x] = data.Registers[x] * 2
}

func Random(data *GameData, x uint16, nn byte) {
	r := byte(rand.Intn(255))
	data.Registers[x] = r & nn
}

func DecimalConversion(data *GameData, x uint16) {
	vx := data.Registers[x]
	data.Memory[data.I] = vx / 100
	data.Memory[data.I+1] = (vx % 100) / 10
	data.Memory[data.I+2] = vx % 10
}

func StoreMemory(data *GameData, x uint16) {
	for i := 0; i <= int(x); i++ {
		data.Memory[data.I+uint16(i)] = data.Registers[i]
	}
}

func LoadMemory(data *GameData, x uint16) {
	for i := 0; i <= int(x); i++ {
		data.Registers[i] = data.Memory[data.I+uint16(i)]
	}
}

func SkipIfKeyPressed(data *GameData, x uint16, key string, isEqualOptional ...bool) {
	isEqual := true
	if len(isEqualOptional) > 0 {
		isEqual = isEqualOptional[0]
	}

	result := data.Registers[x] == helpers.GetKey(key)
	if key != "" && result == isEqual {
		Skip(data)
	}
}

func WaitForKey(data *GameData, x uint16, key string) {
	if key != "" {
		SetRegister(data, x, helpers.GetKey(key))
	} else {
		data.PC -= 2
	}
}
