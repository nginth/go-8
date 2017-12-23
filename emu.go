package main

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
	delay_timer uint8
	sound_timer uint8
	// stack and stack pointer
	stack [16]uint8
	sp    uint8
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
	emu.pc = 0x0000
	memset(emu.gfx[:], 0x00)
	emu.delay_timer = 0x00
	emu.sound_timer = 0x00
	memset(emu.stack[:], 0x00)
	emu.sp = 0x00
	memset(emu.key[:], 0x00)
	emu.drawFlag = 0x00
}

func (emu *Go8) emulateCycle() {
	// fetch opcode
	emu.opcode = emu.getOpcode()
	// decode opcode
	// execute opcode
	// update timers
}

func (emu *Go8) getOpcode() uint16 {
	return uint16(emu.memory[emu.pc])<<8 | uint16(emu.memory[emu.pc+1])
}

func memset(arr []uint8, val uint8) {
	for i := 0; i < len(arr); i++ {
		arr[i] = val
	}
}
