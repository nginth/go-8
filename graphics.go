package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	width       = 640 + 20
	height      = 320 + 20
	thickness   = 10
	pixelSize   = 10
	pixelWidth  = 64
	pixelHeight = 32
)

func setupGraphics() *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:  "GO8",
		Bounds: pixel.R(0, 0, width, height),
	}
	window, err := pixelgl.NewWindow(cfg)
	check(err)
	window.Clear(colornames.Black)
	return window
}

func updateWindow(window *pixelgl.Window, gfx []uint8) {
	window.Clear(colornames.Black)
	drawGfx(window, gfx[:])
	window.Update()
}

func drawGfx(window *pixelgl.Window, gfx []uint8) {
	imd := imdraw.New(nil)
	imd.Color = colornames.White
	for y := 0; y < pixelHeight; y++ {
		for x := 0; x < pixelWidth; x++ {
			if gfx[x+y*pixelWidth] == 1 {
				createSquare(imd, x, pixelHeight-y)
			}
		}
	}
	imd.Draw(window)
}

func createSquare(imd *imdraw.IMDraw, xpos, ypos int) {
	x := float64(pixelSize * xpos)
	y := float64(pixelSize * ypos)
	imd.Push(pixel.V(x, y))           // bottom left
	imd.Push(pixel.V(x+pixelSize, y)) // bottom right
	imd.Rectangle(thickness)
}
