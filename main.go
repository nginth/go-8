package main

import (
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func run() {
	go8 := Go8{}
	go8.initialize()
	go8.loadROM("roms/tetris.ch8")
	timerChan := time.NewTicker(time.Second / 60).C
	cycleChan := time.NewTicker(time.Second / 300).C

	window := setupGraphics()
	for !window.Closed() {
		select {
		case <-cycleChan:
			go8.emulateCycle()
			//fmt.Printf("%x\n", go8.opcode)
			if go8.drawFlag {
				updateWindow(window, go8.gfx[:])
			}
			go8.setKeys(window)
		case <-timerChan:
			go8.updateTimers()
		default:
			// don't block
		}
	}
}

func main() {
	pixelgl.Run(run)
}
