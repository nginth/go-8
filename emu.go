package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"

	"github.com/faiface/pixel/pixelgl"
)

const (
	spriteWidth = 8
	spriteMem   = 0x50
	startPc     = 0x200
)

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
	key      [16]uint8
	drawFlag bool
}

var mathFnTable = []func(*Go8){
	0x0000: (*Go8).setRegs,
	0x0001: (*Go8).orRegs,
	0x0002: (*Go8).andRegs,
	0x0003: (*Go8).xorRegs,
	0x0004: (*Go8).addRegs,
	0x0005: (*Go8).subRegs,
	0x0006: (*Go8).rshift,
	0x0007: (*Go8).subRegsReverse,
	0x000E: (*Go8).lshift,
}

var utilFnTable = []func(*Go8){
	0x0007: (*Go8).storeDelay,
	0x000A: (*Go8).getKey,
	0x0015: (*Go8).setDelay,
	0x0018: (*Go8).setSound,
	0x001E: (*Go8).addToIndex,
	0x0029: (*Go8).getSprite,
	0x0033: (*Go8).storeBCD,
	0x0055: (*Go8).regDump,
	0x0065: (*Go8).regLoad,
}

var fnTable = []func(*Go8){
	0x0000: func(emu *Go8) {
		switch emu.opcode & 0xF000 {
		case 0x0000:
			switch emu.opcode & 0x000F {
			case 0x0000:
				emu.clearScreen()
			case 0x000E:
				emu.ret()
			}
		}
	},
	0x1000: (*Go8).jump,
	0x2000: (*Go8).callSubroutine,
	0x3000: (*Go8).ifEqual,
	0x4000: (*Go8).ifNotEqual,
	0x5000: (*Go8).ifEqualReg,
	0x6000: (*Go8).setConstant,
	0x7000: (*Go8).addConstant,
	0x8000: func(emu *Go8) {
		op := mathFnTable[emu.opcode&0x000F]
		if op == nil {
			fmt.Printf("Unknown opcode: %x\n", emu.opcode)
		} else {
			op(emu)
		}
	},
	0x9000: (*Go8).ifNotEqualReg,
	0xA000: (*Go8).setIndex,
	0xB000: (*Go8).addJump,
	0xC000: (*Go8).rand,
	0xD000: (*Go8).draw,
	0xE000: func(emu *Go8) {
		switch emu.opcode & 0x00FF {
		case 0x009E:
			emu.ifPressed()
		case 0x00A1:
			emu.ifNotPressed()
		}
	},
	0xF000: func(emu *Go8) {
		op := utilFnTable[emu.opcode&0x00FF]
		if op == nil {
			fmt.Printf("Unknown opcode: %x\n", emu.opcode)
		} else {
			op(emu)
		}
	},
}

func (emu *Go8) emulateCycle() {
	emu.opcode = emu.getOpcode()
	op := fnTable[emu.opcode&0xF000]
	if op == nil {
		fmt.Printf("Unknown opcode: %x\n", emu.opcode)
	} else {
		op(emu)
	}
	emu.updateTimers()
}

func (emu *Go8) initialize() {
	emu.opcode = 0x0000
	memset(emu.memory[:], 0x00)
	memset(emu.V[:], 0x00)
	emu.index = 0x0000
	emu.pc = startPc
	memset(emu.gfx[:], 0x00)
	emu.delayTimer = 0x00
	emu.soundTimer = 0x00
	memset16(emu.stack[:], 0x00)
	emu.sp = 0x00
	memset(emu.key[:], 0x00)
	emu.drawFlag = false
	for i := 0; i < 80; i++ {
		emu.memory[spriteMem+i] = fontset[i]
	}
}

func (emu *Go8) loadROM(filename string) {
	data, err := ioutil.ReadFile(filename)
	check(err)
	for i := 0; i < len(data); i++ {
		emu.memory[i+0x200] = data[i]
	}
}

func (emu *Go8) getOpcode() uint16 {
	return uint16(emu.memory[emu.pc])<<8 | uint16(emu.memory[emu.pc+1])
}

func (emu *Go8) setKeys(window *pixelgl.Window) {
	for key := 0; key < len(emu.key); key++ {
		emu.key[key] = 0
		button := keymapping[uint8(key)]
		if window.Pressed(button) {
			emu.key[key] = 1
		}
	}
}

func (emu *Go8) updateTimers() {
	if emu.delayTimer > 0 {
		emu.delayTimer--
	}
	if emu.soundTimer > 0 {
		if emu.soundTimer == 1 {
			playSound()
		}
		emu.soundTimer--
	}
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
	if emu.V[x] > emu.V[y] {
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
	if emu.V[y] > emu.V[x] {
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

func (emu *Go8) addJump() {
	emu.pc = uint16(emu.V[0]) + (emu.opcode & 0x0FFF)
}

func (emu *Go8) rand() {
	x := emu.xreg()
	emu.V[x] = uint8(rand.Intn(256)) & uint8((emu.opcode & 0x00FF))
	emu.pc += 2
}

func (emu *Go8) draw() {
	x := emu.V[emu.xreg()]
	y := emu.V[emu.yreg()]
	height := emu.opcode & 0x000F

	emu.V[0xF] = 0
	var yline uint16
	var xline uint16
	for yline = 0; yline < height; yline++ {
		pixelLine := uint16(emu.memory[emu.index+yline])
		for xline = 0; xline < spriteWidth; xline++ {
			if (pixelLine & (0x80 >> xline)) != 0 {
				pixel := ((uint16(x) + xline + ((uint16(y) + yline) * 64)) % 2048)
				emu.V[0xF] = emu.V[0xF] | emu.gfx[pixel]
				// if emu.gfx[pixel] > 0 {
				// 	fmt.Printf("pixelLine %b\n", pixelLine)
				// 	fmt.Printf("pixel: %b\n", pixel)
				// 	fmt.Printf("x: %d, y: %d, xline: %d, yline: %d, height: %d\n", x, y, xline, yline, height)
				// }
				emu.gfx[pixel] ^= 1
			}
		}
	}
	if emu.V[0xF] > 0 {
		// prevent weird bugs where emu.gfx[n] > 1
		emu.V[0xF] = 1
	}
	emu.drawFlag = true
	emu.pc += 2
}

func (emu *Go8) ifPressed() {
	x := emu.V[emu.xreg()]
	if emu.key[x] == 1 {
		emu.pc += 2
	}
	emu.pc += 2
}

func (emu *Go8) ifNotPressed() {
	x := emu.V[emu.xreg()]
	if emu.key[x] != 1 {
		emu.pc += 2
	}
	emu.pc += 2
}

func (emu *Go8) storeDelay() {
	emu.V[emu.xreg()] = emu.delayTimer
	emu.pc += 2
}

func (emu *Go8) getKey() {
	x := emu.xreg()
	for key := 0; key < len(emu.key); key++ {
		if emu.key[key] == 1 {
			emu.V[x] = uint8(key)
			emu.pc += 2
			break
		}
	}
}

func (emu *Go8) setDelay() {
	emu.delayTimer = emu.V[emu.xreg()]
	emu.pc += 2
}

func (emu *Go8) setSound() {
	emu.soundTimer = emu.V[emu.xreg()]
	emu.pc += 2
}

func (emu *Go8) addToIndex() {
	emu.V[0xF] = 0
	if emu.index+uint16(emu.V[emu.xreg()]) > 0xFFF {
		emu.V[0xF] = 1
	}
	emu.index += uint16(emu.V[emu.xreg()])
	emu.pc += 2
}

func (emu *Go8) getSprite() {
	sprite := uint16(emu.V[emu.xreg()])
	emu.index = 0x50 + sprite*5
	emu.pc += 2
}

func (emu *Go8) storeBCD() {
	x := emu.V[emu.xreg()]
	emu.memory[emu.index] = x / 100
	emu.memory[emu.index+1] = (x / 10) % 10
	emu.memory[emu.index+2] = (x % 100) % 10
	emu.pc += 2
}

func (emu *Go8) regDump() {
	x := emu.xreg()
	var i uint16
	for i = 0; i <= x; i++ {
		emu.memory[emu.index] = emu.V[i]
		emu.index++
	}
	emu.pc += 2
}

func (emu *Go8) regLoad() {
	x := emu.xreg()
	var i uint16
	for i = 0; i <= x; i++ {
		emu.V[i] = emu.memory[emu.index]
		emu.index++
	}
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

var keymapping = map[uint8]pixelgl.Button{
	0x1: pixelgl.Key1,
	0x2: pixelgl.Key2,
	0x3: pixelgl.Key3,
	0xC: pixelgl.Key4,
	0x4: pixelgl.KeyQ,
	0x5: pixelgl.KeyW,
	0x6: pixelgl.KeyE,
	0xD: pixelgl.KeyR,
	0x7: pixelgl.KeyA,
	0x8: pixelgl.KeyS,
	0x9: pixelgl.KeyD,
	0xE: pixelgl.KeyF,
	0xA: pixelgl.KeyZ,
	0x0: pixelgl.KeyX,
	0xB: pixelgl.KeyC,
	0xF: pixelgl.KeyV,
}
