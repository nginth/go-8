package main

import (
	"fmt"
	"testing"
)

func TestFetchOpcode(t *testing.T) {
	go8 := Go8{}
	go8.memory[0] = 0xBE
	go8.memory[1] = 0xEF
	expected := uint16(0xBEEF)
	opcode := go8.getOpcode()
	if opcode != expected {
		t.Errorf("Opcode incorrect, got: %d, want %d.", opcode, expected)
	}
}

func TestInitialize(t *testing.T) {
	go8 := Go8{}
	go8.memory[1] = 0xEF
	go8.opcode = 0xBEEF
	go8.V[3] = 0xFF
	go8.index = 0x12
	go8.pc = 0xF3E4
	go8.gfx[2047] = 0xFF
	go8.delay_timer = 0x23
	go8.sound_timer = 0x34
	go8.stack[14] = 0x0F
	go8.sp = 0x22
	go8.key[0] = 0x11
	go8.drawFlag = 0x1
	go8.initialize()
	if !allFieldsZero(&go8) {
		t.Error("Not initialized to zero.")
		fmt.Println(go8)
	}
}

func allFieldsZero(emu *Go8) bool {
	return emu.opcode == 0 &&
		allArrZero(emu.memory[:]) &&
		allArrZero(emu.V[:]) &&
		emu.index == 0 &&
		emu.pc == 0 &&
		allArrZero(emu.gfx[:]) &&
		emu.delay_timer == 0 &&
		emu.sound_timer == 0 &&
		allArrZero(emu.stack[:]) &&
		emu.sp == 0 &&
		allArrZero(emu.key[:]) &&
		emu.drawFlag == 0
}

func allArrZero(arr []uint8) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] != 0 {
			return false
		}
	}
	return true
}
