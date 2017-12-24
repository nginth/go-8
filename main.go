package main

import (
	"math/rand"

	"github.com/faiface/pixel/pixelgl"
)

func run() {
	// setupInput
	go8 := Go8{}
	// initialize
	go8.initialize()
	// TODO: load ROM
	// go8.loadROM("rom")
	window := setupGraphics()
	// emulation loop
	for !window.Closed() {
		//go8.emulateCycle()
		if go8.drawFlag != 0 {
			// TODO: draw graphics
		}
		for i := 0; i < len(go8.gfx); i++ {
			go8.gfx[i] = uint8(rand.Intn(2))
		}
		updateWindow(window, go8.gfx[:])
		// TODO: store key press state
		// go8.setKeys()
	}
}

func main() {
	pixelgl.Run(run)
}
