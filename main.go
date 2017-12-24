package main

import "github.com/faiface/pixel/pixelgl"

func run() {
	// setupInput
	go8 := Go8{}
	// initialize
	go8.initialize()
	// TODO: load ROM
	// go8.loadROM("rom")
	win := setupGraphics()
	// emulation loop
	for !win.Closed() {
		//go8.emulateCycle()
		if go8.drawFlag != 0 {
			// TODO: draw graphics
		}
		win.Update()
		// TODO: store key press state
		// go8.setKeys()
	}
}

func main() {
	pixelgl.Run(run)
}
