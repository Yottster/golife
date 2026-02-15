package main

import (
	"syscall/js"
	"unsafe"
)

func Render(u *Universe, dst []uint32, colors []uint32) {
	
	for i, cell := range u.cells {
		dst[i] = colors[cell]
	}
}

func main() {
	window := js.Global()
	document := window.Get("document")

	body := document.Get("body")
	bodyStyle := body.Get("style")

	deadColor := Color{0,0,0,200}
	liveColor := Color{0,0,0xFF,0xFF}

	colors := []uint32{ deadColor.Packed(), liveColor.Packed() }

	bodyStyle.Set("background", deadColor.Solid().String())
	
	innerWidth := window.Get("innerWidth").Int()
	innerHeight := window.Get("innerHeight").Int()

	cellSize := window.Get("cellSize").Int()
	width := innerWidth/cellSize
	height := innerHeight/cellSize
	
	pixels := make([]byte, 4 * width * height)
	ptr := unsafe.Pointer(&pixels[0])
	u32Buffer := 
		unsafe.Slice((*uint32)(ptr), width * height)

	nextUniverse := NewUniverse(width, height)
	currentUniverse := NewUniverse(width, height)
	currentUniverse.seed()

	performance := window.Get("performance")
	frameTimer := NewFrameTimer(performance, 30)

	var renderFrame js.Func
	var status string
	
	renderFrame = js.FuncOf(func (this js.Value, args []js.Value) interface {} {

		start := frameTimer.Start()

		Render(currentUniverse, u32Buffer, colors)

		window.Call(
			"updateCanvas", 
			uintptr(ptr), 
			width, 
			height,
			status)

		currentUniverse.Next(nextUniverse)
		currentUniverse, nextUniverse = nextUniverse, currentUniverse

		update := frameTimer.End(start)
		if update {
			status = frameTimer.Status()
		}

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	
	js.Global().Call("requestAnimationFrame", renderFrame)
	select {}
}
