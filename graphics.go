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

var keymapping = map[uint8]pixelgl.Button{
	0x1: pixelgl.Key1,
	0x2: pixelgl.Key2,
	0x3: pixelgl.Key3,
	0xC: pixelgl.Key4,
	0x4: pixelgl.KeyQ,
	0x5: pixelgl.KeyW,
	0x6: pixelgl.KeyE,
	0xD: pixelgl.KeyR,
	0x7: pixelgl.KeyA,
	0x8: pixelgl.KeyS,
	0x9: pixelgl.KeyD,
	0xE: pixelgl.KeyF,
	0xA: pixelgl.KeyZ,
	0x0: pixelgl.KeyX,
	0xB: pixelgl.KeyC,
	0xF: pixelgl.KeyV,
}

// GraphicsDevice - a generic graphics device interface
type GraphicsDevice interface {
	updateWindow(gfx []uint8)
	closed() bool
	pressed(key int) bool
}

// Graphics - a pixel implementation of GraphicsDevice
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
				graphics.createPixel(imd, x, pixelHeight-y)
			}
		}
	}
	imd.Draw(graphics.window)
}

func (graphics *Graphics) closed() bool {
	return graphics.window.Closed()
}

func (graphics *Graphics) pressed(button int) bool {
	return graphics.window.Pressed(keymapping[uint8(button)])
}

func (graphics *Graphics) createPixel(imd *imdraw.IMDraw, xpos, ypos int) {
	x := float64(pixelSize * xpos)
	y := float64(pixelSize * ypos)
	imd.Push(pixel.V(x, y))           // bottom left
	imd.Push(pixel.V(x+pixelSize, y)) // bottom right
	imd.Rectangle(thickness)
}
