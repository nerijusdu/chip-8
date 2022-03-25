package main

import (
	"chip8/helpers"
	"chip8/instructions"
	. "chip8/models"
	"encoding/binary"
	"os"
	"time"
)

var font []byte = []byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

// Keypad
// 1 	2 	3 	C
// 4 	5 	6 	D
// 7 	8 	9 	E
// A 	0 	B 	F

const pixelCount = 64 * 32
const clockSpeed = 700

var data *GameData

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Init() *GameData {
	data = &GameData{
		Pixels:     make(Pixels, pixelCount),
		Stack:      make(Stack, 100),
		Registers:  make([]byte, 16),
		Memory:     make([]byte, 4096),
		DelayTimer: 255,
		SoundTimer: 255,
		PC:         0x200,
		I:          0,
		PixelCount: pixelCount,
		ClockSpeed: clockSpeed,
	}

	for i, v := range font {
		data.Memory[i] = v
	}

	program := ReadFile()
	for i, v := range program {
		data.Memory[0x200+i] = v
	}

	return data
}

func ReadFile() []byte {
	dat, err := os.ReadFile("roms/IBM Logo.ch8")
	check(err)

	return dat
}

func GameLoop(refresh func()) {
	cycles := 0
	for {
		opcode := binary.BigEndian.Uint16([]byte{data.Memory[data.PC], data.Memory[data.PC+1]})
		data.PC += 2
		Execute(opcode)

		data.UpdateTimers()

		refresh()
		cycles++
		time.Sleep(100 * time.Millisecond)
	}
}

func Execute(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	nnn := opcode & 0x0FFF
	nn := helpers.GetByteFrom16(opcode & 0x00FF)
	n := opcode & 0x000F

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			instructions.ClearScreen(data)
			break
		case 0x00EE:
			break
		}

		break
	case 0x1000:
		instructions.Jump(data, nnn)
		break
	case 0x2000:
		break
	case 0x3000:
		break
	case 0x4000:
		break
	case 0x5000:
		break
	case 0x6000:
		instructions.SetRegister(data, x, nn)
		break
	case 0x7000:
		instructions.AddValueToRegister(data, x, nn)
		break
	case 0x8000:
		switch opcode & 0xF {
		case 0x0:
			break
		case 0x1:
			break
		case 0x2:
			break
		case 0x3:
			break
		case 0x4:
			break
		case 0x5:
			break
		case 0x6:
			break
		case 0x7:
			break
		case 0xE:
			break
		}

		break
	case 0x9000:
		break
	case 0xA000:
		instructions.SetIndexRegister(data, nnn)
		break
	case 0xB000:
		break
	case 0xC000:
		break
	case 0xD000:
		instructions.DisplayDraw(data, x, y, n)
		break
	case 0xE000:
		switch opcode & 0xFF {
		case 0x9E:
			break
		case 0xA1:
			break
		}

		break
	case 0xF000:
		switch opcode & 0xFF {
		case 0x07:
			break
		case 0x0A:
			break
		case 0x15:
			break
		case 0x18:
			break
		case 0x1E:
			break
		case 0x29:
			break
		case 0x33:
			break
		case 0x55:
			break
		case 0x65:
			break
		}

		break

	}
}
