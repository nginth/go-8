package main

import "github.com/go-gl/glfw/v3.2/glfw"

const (
	width  = 640
	height = 320
)

func setupGraphics() *glfw.Window {
	err := glfw.Init()
	check(err)

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "GO-8", nil, nil)
	check(err)

	window.MakeContextCurrent()
	return window
}

func terminateGraphics() {
	glfw.Terminate()
}

func updateWindow(window *glfw.Window) {
	window.SwapBuffers()
	glfw.PollEvents()
}

func render(window *glfw.Window, gfx []uint8) {

}
