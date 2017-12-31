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

type Graphics struct {
	window *pixelgl.Window
}

func newGraphics() *Graphics {
	cfg := pixelgl.WindowConfig{
		Title:  "GO8",
		Bounds: pixel.R(0, 0, width, height),
	}
	window, err := pixelgl.NewWindow(cfg)
	check(err)
	window.Clear(colornames.Black)
	return &Graphics{window: window}
}

func (graphics *Graphics) updateWindow(gfx []uint8) {
	graphics.window.Clear(colornames.Black)
	graphics.drawGfx(gfx[:])
	graphics.window.Update()
}

func (graphics *Graphics) drawGfx(gfx []uint8) {
	imd := imdraw.New(nil)
	imd.Color = colornames.White
	for y := 0; y < pixelHeight; y++ {
		for x := 0; x < pixelWidth; x++ {
			if gfx[x+y*pixelWidth] == 1 {
				graphics.createSquare(imd, x, pixelHeight-y)
			}
		}
	}
	imd.Draw(graphics.window)
}

func (graphics *Graphics) createSquare(imd *imdraw.IMDraw, xpos, ypos int) {
	x := float64(pixelSize * xpos)
	y := float64(pixelSize * ypos)
	imd.Push(pixel.V(x, y))           // bottom left
	imd.Push(pixel.V(x+pixelSize, y)) // bottom right
	imd.Rectangle(thickness)
}
