package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	go8.delayTimer = 0x23
	go8.soundTimer = 0x34
	go8.stack[14] = 0x0F
	go8.sp = 0x22
	go8.key[0] = 0x11
	go8.drawFlag = true
	go8.initialize()
	if !allFieldsInit(&go8) {
		t.Error("Not initialized to zero.")
		fmt.Println(go8)
	}
}

func TestReadROM(t *testing.T) {
	ex, err := os.Executable()
	check(err)
	path := filepath.Dir(ex)
	f, err := ioutil.TempFile(path, "")
	tmprom := f.Name()
	check(err)

	buf := []byte{0x55, 0x55, 0x55, 0x55}
	f.Write(buf)
	f.Close()

	go8 := Go8{}
	go8.loadROM(tmprom)
	for i := 0; i < len(buf); i++ {
		if go8.memory[i+512] != 0x55 {
			t.Errorf("Invalid memory state. Got %d, wanted %d", go8.memory[i], 0x55)
		}
	}
	os.Remove(tmprom)
}

func TestEmulateCycleSubroutine(t *testing.T) {
	go8 := Go8{}
	go8.initialize()
	go8.pc = 0x512
	go8.memory[0x512] = 0x22
	go8.memory[0x513] = 0x22
	go8.delayTimer = 2
	go8.emulateCycle()
	if go8.stack[0] != 0x512 {
		t.Errorf("Wrong stack value. Got %x, expected %x.", go8.stack[0], 0x512)
	}
	if go8.sp != 0x1 {
		t.Errorf("Wrong sp. Got %x, expected %x.", go8.sp, 0x1)
	}
	if go8.delayTimer != 1 {
		t.Errorf("Wrong delay timer. Got %d, expected %d.", go8.delayTimer, 1)
	}
	checkPc(0x222, go8.pc, t)
}

func TestCallSubroutine(t *testing.T) {
	go8 := Go8{}
	go8.initialize()
	go8.opcode = 0x2222
	go8.pc = 0x512
	go8.callSubroutine()
	if go8.stack[0] != 0x512 {
		t.Errorf("Wrong stack value. Got %x, expected %x.", go8.stack[0], 0x512)
	}
	if go8.sp != 0x1 {
		t.Errorf("Wrong sp. Got %x, expected %x.", go8.sp, 0x1)
	}
	checkPc(0x222, go8.pc, t)
}

func TestReturnSubroutine(t *testing.T) {
	go8 := Go8{}
	go8.initialize()
	go8.opcode = 0x00EE
	go8.stack[0] = 0x512
	go8.sp = 0x1
	go8.ret()
	if go8.sp != 0x0 {
		t.Errorf("Wrong sp. Got %x, expected %x.", go8.sp, 0x1)
	}
	checkPc(0x512+2, go8.pc, t)
}

func TestJump(t *testing.T) {
	go8 := Go8{}
	go8.initialize()
	go8.opcode = 0x1111
	go8.pc = 0x512
	go8.jump()
	checkPc(0x111, go8.pc, t)
}

func TestClearScreen(t *testing.T) {
	go8 := Go8{}
	go8.initialize()
	go8.opcode = 0x1111
	go8.pc = 0x512
	for i := 0; i < len(go8.gfx); i++ {
		go8.gfx[i] = 0x12
	}
	go8.clearScreen()
	if !allArrZero(go8.gfx[:]) {
		t.Errorf("Gfx not cleared. Got %v, expected all zeroes.", go8.gfx[:])
	}
	checkPc(0x512+2, go8.pc, t)
}

func TestAddRegs(t *testing.T) {
	go8 := Go8{}
	go8.initialize()
	go8.opcode = 0x8124
	go8.V[1] = 2
	go8.V[2] = 2
	go8.pc = 0x512
	go8.addRegs()
	if go8.V[1] != 4 {
		t.Errorf("Wrong value in V[1]. Got %d, expected %d", go8.V[1], 4)
	}
	checkPc(0x512+2, go8.pc, t)
	go8.V[1] = 255
	go8.V[2] = 1
	go8.addRegs()
	if go8.V[1] != 0 {
		t.Errorf("Wrong value in V[1]. Got %d, expected %d", go8.V[1], 0)
	}
	if go8.V[0xF] != 1 {
		t.Errorf("Carry flag not set. Got %d, expected %d", go8.V[0xF], 1)
	}
}

func allFieldsInit(emu *Go8) bool {
	return emu.opcode == 0 &&
		allArrZero(emu.memory[80:]) && // fontset stored < 0x50
		allArrZero(emu.V[:]) &&
		emu.index == 0 &&
		emu.pc == 0x0200 &&
		allArrZero(emu.gfx[:]) &&
		emu.delayTimer == 0 &&
		emu.soundTimer == 0 &&
		allArrZero16(emu.stack[:]) &&
		emu.sp == 0 &&
		allArrZero(emu.key[:]) &&
		emu.drawFlag == false
}

func checkPc(expected int, actual uint16, t *testing.T) {
	if actual != uint16(expected) {
		t.Errorf("Wrong pc. Got %x, expected %x.", actual, expected)
	}
}

func allArrZero(arr []uint8) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] != 0 {
			return false
		}
	}
	return true
}

func allArrZero16(arr []uint16) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] != 0 {
			return false
		}
	}
	return true
}
