package main

import (
	"chip8/helpers"
	"chip8/instructions"
	. "chip8/models"
	"encoding/binary"
	"fmt"
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

// TODO: investigate if font is rendered correctly

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
		Stack:      *NewStack(0),
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
	dat, err := os.ReadFile("roms/keypad_test.ch8")
	check(err)

	return dat
}

func GameLoop(refresh func(), messages <-chan string) {
	helpers.Schedule(time.Duration(500)*time.Millisecond, refresh)
	helpers.Schedule(time.Duration(1000000/60)*time.Microsecond, data.UpdateTimers) // 60hz
	waitTime := time.Duration(1000000/clockSpeed) * time.Microsecond                // 700hz
	cycles := 0

	for {
		key := ""
		select {
		case key = <-messages:
			fmt.Printf("-------Key is pressed %s\n", key)
			break
		default:
			key = ""
			break
		}

		opcode := binary.BigEndian.Uint16([]byte{data.Memory[data.PC], data.Memory[data.PC+1]})
		// fmt.Println("Key:", key, data.PC, fmt.Sprintf("%X", opcode)) // TODO: keypad only registers first key
		data.PC += 2
		Execute(opcode, key)

		cycles++
		time.Sleep(waitTime)
	}
}

func Execute(opcode uint16, key string) {
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
			instructions.ReturnFromSubroutine(data)
			break
		}

		break
	case 0x1000:
		instructions.Jump(data, nnn)
		break
	case 0x2000:
		instructions.StartSubroutine(data, nnn)
		break
	case 0x3000:
		instructions.SkipIfEqual(data, x, nn)
		break
	case 0x4000:
		instructions.SkipIfEqual(data, x, nn, false)
		break
	case 0x5000:
		instructions.SkipIfEqualRegisters(data, x, y)
		break
	case 0x6000:
		instructions.SetRegister(data, x, nn)
		break
	case 0x7000:
		instructions.AddValueToRegister(data, x, nn)
		break
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			data.Registers[x] = data.Registers[y]
			break
		case 0x0001:
			data.Registers[x] = data.Registers[x] | data.Registers[y]
			break
		case 0x0002:
			data.Registers[x] = data.Registers[x] & data.Registers[y]
			break
		case 0x0003:
			data.Registers[x] = data.Registers[x] ^ data.Registers[y]
			break
		case 0x0004:
			instructions.AddWithCarry(data, x, y)
			break
		case 0x0005:
			instructions.SubtractWithCarry(data, x, x, y)
			break
		case 0x0006:
			instructions.ShiftRight(data, x)
			break
		case 0x0007:
			instructions.SubtractWithCarry(data, x, y, x)
			break
		case 0x000E:
			instructions.ShiftLeft(data, x)
			break
		}

		break
	case 0x9000:
		instructions.SkipIfEqualRegisters(data, x, y, false)
		break
	case 0xA000:
		instructions.SetIndexRegister(data, nnn)
		break
	case 0xB000:
		instructions.JumpWithOffset(data, nnn)
		break
	case 0xC000:
		instructions.Random(data, x, nn)
		break
	case 0xD000:
		instructions.DisplayDraw(data, x, y, n)
		break
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			instructions.SkipIfKeyPressed(data, x, key)
			break
		case 0x00A1:
			instructions.SkipIfKeyPressed(data, x, key, false)
			break
		}

		break
	case 0xF000:
		switch opcode & 0xFF {
		case 0x0007:
			data.Registers[x] = data.DelayTimer
			break
		case 0x000A:
			instructions.WaitForKey(data, x, key)
			break
		case 0x0015:
			data.DelayTimer = data.Registers[x]
			break
		case 0x0018:
			data.SoundTimer = data.Registers[x]
			break
		case 0x001E:
			data.I += uint16(data.Registers[x])
			break
		case 0x0029:
			data.I = uint16(data.Registers[x]) * uint16(0x05)
			break
		case 0x0033:
			instructions.DecimalConversion(data, x)
			break
		case 0x0055:
			instructions.StoreMemory(data, x)
			break
		case 0x0065:
			instructions.LoadMemory(data, x)
			break
		}

		break

	}
}
