package models

type GameData struct {
	Pixels     Pixels
	Stack      Stack
	Registers  []byte
	Memory     []byte
	DelayTimer uint8
	SoundTimer uint8
	PC         uint16
	I          uint16
	PixelCount int
	ClockSpeed int
}

func (gameData *GameData) UpdateTimers() {
	if gameData.DelayTimer > 0 {
		gameData.DelayTimer -= 0
	}
	if gameData.SoundTimer > 0 {
		gameData.SoundTimer -= 0
	}
}

type Stack []uint16

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(data uint16) {
	*s = append(*s, data)
}

func (s *Stack) Pop() uint16 {
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element
}

type Pixels []bool

func (p *Pixels) Get(x uint16, y uint16) bool {
	pos := x + y*64
	if x > 63 || y > 31 {
		return false
	}
	return (*p)[pos]
}

func (p *Pixels) Set(x uint16, y uint16, value bool) {
	pos := x + y*64
	if x > 63 || y > 31 {
		return
	}

	(*p)[pos] = value
}
