package main

import (
	"flag"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func run() {
	rom, timerFreq, clockFreq := getFlags()
	go8 := Go8{}
	go8.initialize()
	go8.loadROM(rom)
	timerChan := time.NewTicker(timerFreq).C
	cycleChan := time.NewTicker(clockFreq).C

	window := setupGraphics()
	for !window.Closed() {
		select {
		case <-cycleChan:
			go8.emulateCycle()
			if go8.drawFlag {
				updateWindow(window, go8.gfx[:])
				go8.drawFlag = false
			}
			go8.setKeys(window)
		case <-timerChan:
			go8.updateTimers()
		}
	}
}

func getFlags() (string, time.Duration, time.Duration) {
	rom := flag.String("rom", "roms/tetris.ch8", "Path to rom.")
	timerFreq := flag.Int("timerFreq", 60, "Timer frequency in Hz.")
	clockFreq := flag.Int("clockFreq", 300, "Clock speed in Hz.")
	flag.Parse()
	t := time.Duration((int(time.Second) / *timerFreq))
	c := time.Duration((int(time.Second) / *clockFreq))
	return *rom, t, c
}

func main() {
	pixelgl.Run(run)
}
