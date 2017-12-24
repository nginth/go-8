package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	width     = 640
	height    = 320
	pixelSize = 10
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

func updateWindow(window *pixelgl.Window) {
	window.Clear(colornames.Black)
	createSquare(1, 0).Draw(window)
	window.Update()
}

func createSquare(x, y float64) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(1, 1, 1)
	x *= pixelSize
	y *= pixelSize
	imd.Push(pixel.V(x, y))                     // bottom left
	imd.Push(pixel.V(x+pixelSize, y))           // bottom right
	imd.Push(pixel.V(x, y+pixelSize))           // top left
	imd.Push(pixel.V(x+pixelSize, y+pixelSize)) // top right
	imd.Rectangle(0)
	return imd
}
