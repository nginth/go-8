package main

import (
	"flag"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func run() {
	timerFreq, clockFreq := getFlags()
	go8 := Go8{}
	go8.initialize()
	go8.loadROM("roms/tetris.ch8")
	timerChan := time.NewTicker(timerFreq).C
	cycleChan := time.NewTicker(clockFreq).C

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

func getFlags() (time.Duration, time.Duration) {
	timerFreq := flag.Int("timerFreq", 60, "Timer frequency in Hz.")
	clockFreq := flag.Int("clockFreq", 300, "Clock speed in Hz.")
	flag.Parse()
	t := time.Duration((int(time.Second) / *timerFreq))
	c := time.Duration((int(time.Second) / *clockFreq))
	return t, c
}

func main() {
	pixelgl.Run(run)
}
