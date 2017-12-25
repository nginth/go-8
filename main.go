package main

import (
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	// setupInput
	go8 := Go8{}
	// initialize
	go8.initialize()
	go8.loadROM("roms/tetris.ch8")
	window := setupGraphics()
	go8.input = window
	// emulation loop
	for !window.Closed() {
		go8.emulateCycle()
		//fmt.Printf("%x\n", go8.opcode)
		if go8.drawFlag {
			updateWindow(window, go8.gfx[:])
		}
		go8.setKeys(window)
		// go about 60Hz
		//time.Sleep(time.Duration(17) * time.Millisecond)
	}
}

func main() {
	pixelgl.Run(run)
}
