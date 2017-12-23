package main

var opcode uint16
var memory [4096]uint8

// all registers V0-VF
var V [16]uint8

// index and program counter registers
var index uint16
var pc uint16

// 64 x 32 px screen, black or white
var gfx [64 * 32]uint8

// timers
var delay_timer uint8
var sound_timer uint8

var stack [16]uint8
var sp uint8

// keypad (input device)
var key [16]uint8

func main() {
	// setupGraphics
	// setupInput

	// initialize
	initialize()
	// load ROM

	// emulation loop
}

func initialize() {
	opcode = 0x0000
	memset(memory[:], 0x00)
	memset(V[:], 0x00)
	index = 0x0000
	pc = 0x0000
	memset(gfx[:], 0x00)
	delay_timer = 0x00
	sound_timer = 0x00
	memset(stack[:], 0x00)
	sp = 0x00
	memset(key[:], 0x00)
}

func memset(arr []uint8, val uint8) {
	for i := 0; i < len(arr); i++ {
		arr[i] = val
	}
}
