package main

import (
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
		go8.emulateCycle()
		if go8.drawFlag {
			updateWindow(window, go8.gfx[:])
		}
		go8.setKeys(window)
	}
}

func main() {
	pixelgl.Run(run)
}
