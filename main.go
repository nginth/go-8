package main

import "runtime"

func init() {
	runtime.LockOSThread()
}

func main() {
	window := setupGraphics()
	defer terminateGraphics()
	// setupInput
	go8 := Go8{}
	// initialize
	go8.initialize()
	// TODO: load ROM
	// go8.loadROM("rom")

	// emulation loop
	for !window.ShouldClose() {
		go8.emulateCycle()
		updateWindow(window)
		if go8.drawFlag != 0 {
			// TODO: draw graphics
		}
		// TODO: store key press state
		// go8.setKeys()
	}
}
