package main

import (
	"flag"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func run() {
	rom, timerFreq, clockFreq := getFlags()
	go8 := newGo8()
	go8.loadROM(rom)
	timerChan := time.NewTicker(timerFreq).C
	cycleChan := time.NewTicker(clockFreq).C

	for !go8.graphics.window.Closed() {
		select {
		case <-cycleChan:
			go8.emulateCycle()
			if go8.drawFlag {
				go8.updateWindow()
				go8.drawFlag = false
			}
			go8.setKeys()
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
