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

// Go8 - CHIP-8 emulator
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
	drawFlag bool
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
	emu.drawFlag = false
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
	case 0x0000:
		switch emu.opcode & 0x000F {
		case 0x0000:
			emu.clearScreen()
		case 0x000E:
			emu.ret()
		}
	case 0x1000:
		emu.jump()
	case 0x2000:
		emu.callSubroutine()
	case 0x3000:
		emu.ifEqual()
	case 0x4000:
		emu.ifNotEqual()
	case 0x5000:
		emu.ifEqualReg()
	case 0x6000:
		emu.setConstant()
	case 0x7000:
		emu.addConstant()
	case 0x8000:
		switch emu.opcode & 0x000F {
		case 0x0000:
			emu.setRegs()
		case 0x0001:
			emu.orRegs()
		case 0x0002:
			emu.andRegs()
		case 0x0003:
			emu.xorRegs()
		case 0x0004:
			emu.addRegs()
		case 0x0005:
			emu.subRegs()
		case 0x0006:
			emu.rshift()
		case 0x0007:
			emu.subRegsReverse()
		case 0x000E:
			emu.lshift()
		}
	case 0x9000:
		emu.ifNotEqualReg()
	case 0xA000:
		emu.setIndex()
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

func (emu *Go8) ret() {
	emu.pc = emu.stack[emu.sp-1] + 2
	emu.sp--
}

func (emu *Go8) clearScreen() {
	memset(emu.gfx[:], 0)
	emu.drawFlag = true
	emu.pc += 2
}

func (emu *Go8) ifEqual() {
	x := emu.xreg()
	n := emu.opcode & 0x00FF
	if emu.V[x] == uint8(n) {
		emu.pc += 2
	}
	emu.pc += 2
}

func (emu *Go8) ifNotEqual() {
	x := emu.xreg()
	n := emu.opcode & 0x00FF
	if emu.V[x] != uint8(n) {
		emu.pc += 2
	}
	emu.pc += 2
}

func (emu *Go8) ifEqualReg() {
	x := emu.xreg()
	y := emu.yreg()
	if emu.V[x] == emu.V[y] {
		emu.pc += 2
	}
	emu.pc += 2
}

func (emu *Go8) ifNotEqualReg() {
	x := emu.xreg()
	y := emu.yreg()
	if emu.V[x] != emu.V[y] {
		emu.pc += 2
	}
	emu.pc += 2
}

func (emu *Go8) setConstant() {
	x := emu.xreg()
	n := emu.opcode & 0x00FF
	emu.V[x] = uint8(n)
	emu.pc += 2
}

func (emu *Go8) addConstant() {
	x := emu.xreg()
	n := emu.opcode & 0x00FF
	// by spec, carry flag is not changed on constant add
	emu.V[x] += uint8(n)
	emu.pc += 2
}

func (emu *Go8) setRegs() {
	x := emu.xreg()
	y := emu.yreg()
	emu.V[x] = emu.V[y]
	emu.pc += 2
}

func (emu *Go8) orRegs() {
	x := emu.xreg()
	y := emu.yreg()
	emu.V[x] |= emu.V[y]
	emu.pc += 2
}

func (emu *Go8) andRegs() {
	x := emu.xreg()
	y := emu.yreg()
	emu.V[x] &= emu.V[y]
	emu.pc += 2
}

func (emu *Go8) xorRegs() {
	x := emu.xreg()
	y := emu.yreg()
	emu.V[x] ^= emu.V[y]
	emu.pc += 2
}

func (emu *Go8) addRegs() {
	x := emu.xreg()
	y := emu.yreg()
	if emu.V[y] > 0xFF-emu.V[x] {
		emu.V[0xF] = 1
	} else {
		emu.V[0xF] = 0
	}
	emu.V[x] += emu.V[y]
	emu.pc += 2
}

func (emu *Go8) subRegs() {
	x := emu.xreg()
	y := emu.yreg()
	if emu.V[y] > emu.V[x] {
		emu.V[0xF] = 1
	} else {
		emu.V[0xF] = 0
	}
	emu.V[x] -= emu.V[y]
	emu.pc += 2
}

func (emu *Go8) subRegsReverse() {
	x := emu.xreg()
	y := emu.yreg()
	if emu.V[x] > emu.V[y] {
		emu.V[0xF] = 1
	} else {
		emu.V[0xF] = 0
	}
	emu.V[x] = emu.V[y] - emu.V[x]
	emu.pc += 2
}

func (emu *Go8) rshift() {
	x := emu.xreg()
	y := emu.yreg()
	emu.V[0xF] = emu.V[y] & 0x01
	emu.V[y] >>= 1
	emu.V[x] = emu.V[y]
	emu.pc += 2
}

func (emu *Go8) lshift() {
	x := emu.xreg()
	y := emu.yreg()
	emu.V[0xF] = (emu.V[y] & 0x80) >> 7
	emu.V[y] <<= 1
	emu.V[x] = emu.V[y]
	emu.pc += 2
}

func (emu *Go8) setIndex() {
	emu.index = emu.opcode & 0x0FFF
	emu.pc += 2
}

func (emu *Go8) xreg() uint16 {
	return emu.opcode & 0x0F00 >> 8
}

func (emu *Go8) yreg() uint16 {
	return emu.opcode & 0x00F0 >> 4
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
