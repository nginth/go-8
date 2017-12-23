package main

func main() {
	setupGraphics()
	// setupInput
	go8 := Go8{}
	// initialize
	go8.V[0] = 12
	go8.initialize()
	// TODO: load ROM
	// go8.loadROM("rom")

	// emulation loop
	for {
		go8.emulateCycle()
		if go8.drawFlag != 0 {
			// TODO: draw graphics
		}
		// TODO: store key press state
		// go8.setKeys()
	}
}
