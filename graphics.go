package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width  = 640
	height = 320
)

func setupGraphics() *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:  "GO8",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	window, err := pixelgl.NewWindow(cfg)
	check(err)
	return window
}
