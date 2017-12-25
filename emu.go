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
	input    *pixelgl.Window
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
		emu.memory[spriteMem+i] = fontset[i]
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
	case 0xB000:
		emu.addJump()
	case 0xC000:
		emu.rand()
	case 0xD000:
		emu.draw()
	case 0xE000:
		switch emu.opcode & 0x00FF {
		case 0x009E:
			emu.ifPressed()
		case 0x00A1:
			emu.ifNotPressed()
		}
	case 0xF000:
		switch emu.opcode & 0x00FF {
		case 0x0007:
			emu.storeDelay()
		case 0x000A:
			emu.getKey()
		case 0x0015:
			emu.setDelay()
		case 0x0018:
			emu.setSound()
		case 0x001E:
			emu.addToIndex()
		}
	default:
		fmt.Printf("Unknown opcode: %x\n", emu.opcode)
	}
	// update timers
	emu.updateTimers()
}

func (emu *Go8) getOpcode() uint16 {
	return uint16(emu.memory[emu.pc])<<8 | uint16(emu.memory[emu.pc+1])
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

func (emu *Go8) setKeys(window *pixelgl.Window) {
	for key := 0; key < len(emu.key); key++ {
		emu.key[key] = 0
		button := keymapping[uint8(key)]
		if window.Pressed(button) || window.JustPressed(button) || window.JustReleased(button) {
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
			fmt.Println("BEEP!!")
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

func (emu *Go8) addJump() {
	emu.pc = uint16(emu.V[0]) + (emu.opcode & 0x0FFF)
}

func (emu *Go8) rand() {
	x := emu.xreg()
	emu.V[x] = uint8(rand.Intn(256)) & uint8((emu.opcode & 0x00FF))
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
				pixel := uint16(x) + xline + ((uint16(y) + yline) * 64)
				if emu.gfx[pixel] == 1 {
					emu.V[0xF] = 1
				}
				emu.gfx[pixel] ^= 1
			}
		}
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
	for !emu.input.Closed() {
		for key := 0; key < len(emu.key); key++ {
			button := keymapping[uint8(key)]
			if emu.input.Pressed(button) || emu.input.JustPressed(button) || emu.input.JustReleased(button) {
				emu.V[x] = uint8(key)
				fmt.Println("asdf")
				return
			}
		}
		emu.input.Update()
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
	emu.index += uint16(emu.V[emu.xreg()])
	emu.pc += 2
}

func (emu *Go8) getSprite() {
	sprite := uint16(emu.V[emu.xreg()])
	emu.index = 0x50 + sprite*5
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
