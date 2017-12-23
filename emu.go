package main

import (
	"fmt"
	"io/ioutil"
)

var fontset = [80]uint8{
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

type Go8 struct {
	opcode uint16
	memory [4096]uint8
	// all registers V0-VF
	V [16]uint8
	// index and program counter registers
	index uint16
	pc    uint16
	// 64 x 32 px screen, black or white
	gfx [64 * 32]uint8
	// timers
	delayTimer uint8
	soundTimer uint8
	// stack and stack pointer
	stack [16]uint16
	sp    uint16
	// keypad (input device)
	key [16]uint8
	// graphics
	drawFlag uint8
}

func (emu *Go8) initialize() {
	emu.opcode = 0x0000
	memset(emu.memory[:], 0x00)
	memset(emu.V[:], 0x00)
	emu.index = 0x0000
	emu.pc = 0x0200
	memset(emu.gfx[:], 0x00)
	emu.delayTimer = 0x00
	emu.soundTimer = 0x00
	memset16(emu.stack[:], 0x00)
	emu.sp = 0x00
	memset(emu.key[:], 0x00)
	emu.drawFlag = 0x00
	// load fontset
	for i := 0; i < 80; i++ {
		emu.memory[i] = fontset[i]
	}
}

func (emu *Go8) loadROM(filename string) {
	data, err := ioutil.ReadFile(filename)
	check(err)
	for i := 0; i < len(data); i++ {
		emu.memory[i+512] = data[i]
	}
}

func (emu *Go8) emulateCycle() {
	emu.opcode = emu.getOpcode()
	// execute opcode
	switch emu.opcode & 0xF000 {
	case 0xA000:
		emu.index = emu.opcode & 0x0FFF
		emu.pc += 2
	case 0x1000:
		// opcode 0x1NNN : jump to address NNN
		emu.jump()
	case 0x2000:
		// opcode 0x2NNN : call subroutine at address NNN
		emu.callSubroutine()
	default:
		fmt.Printf("Unknown opcode: %x\n", emu.opcode)
	}
	// update timers
	if emu.delayTimer > 0 {
		emu.delayTimer--
	}
	if emu.soundTimer > 0 {
		if emu.soundTimer == 1 {
			fmt.Println("BEEP!!")
		}
		emu.soundTimer--
	}
}

func (emu *Go8) getOpcode() uint16 {
	return uint16(emu.memory[emu.pc])<<8 | uint16(emu.memory[emu.pc+1])
}

func (emu *Go8) callSubroutine() {
	emu.stack[emu.sp] = emu.pc
	emu.sp++
	emu.pc = emu.opcode & 0x0FFF
}

func (emu *Go8) jump() {
	emu.pc = emu.opcode & 0x0FFF
}

func memset(arr []uint8, val uint8) {
	for i := 0; i < len(arr); i++ {
		arr[i] = val
	}
}

func memset16(arr []uint16, val uint16) {
	for i := 0; i < len(arr); i++ {
		arr[i] = val
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
