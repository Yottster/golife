package main

import (
	"syscall/js"
)

var privateUniverse *Universe

func InitGame(this js.Value, args []js.Value) interface{} {
	width := args[0].Int()
	height := args[1].Int()

	privateUniverse = NewUniverse(width, height)
	privateUniverse.seed()
	
	return nil
}
func Paint(ctx js.Value, 
	universe *Universe, 
	cellSize int,
	status string) {
	
	ctx.Call("clearRect", 0, 0, 
		universe.width * cellSize, 
		universe.height * cellSize)

	ctx.Set("fillStyle", "red")
	ctx.Call("fillText", status, 10, 30)
	
	ctx.Set("fillStyle", "#000000")

	for y := 0; y < universe.height; y++ {
		for x := 0; x < universe.width; x++ {
			if universe.Alive(x, y) {
				ctx.Call("fillRect", 
					x * cellSize, 
					y * cellSize, 
					cellSize, 
					cellSize)
			}
		}
	}
}

func PrintStatus(ctx js.Value, status string) {
	ctx.Set("fillStyle", "red")
	ctx.Call("fillText", status, 10, 30)
}
func main() {
	window := js.Global()
	document := window.Get("document")
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")
	cellSize := 3
	
	innerWidth := window.Get("innerWidth").Int()
	innerHeight := window.Get("innerHeight").Int()

	canvas.Set("width", innerWidth)
	canvas.Set("height", innerHeight)

	width := innerWidth/cellSize
	height := innerHeight/cellSize
	
	currentUniverse := NewUniverse(width, height)
	nextUniverse := NewUniverse(width, height)
	
	currentUniverse.seed()

	performance := window.Get("performance")
	frameTimer := NewFrameTimer(performance, 30)
	ctx.Set("font", "40px monospace")
	var renderFrame js.Func
	var status string
	
	renderFrame = js.FuncOf(func (this js.Value, args []js.Value) interface {} {

		start := frameTimer.Start()

		Paint(ctx, currentUniverse, cellSize, status)
		
		currentUniverse.Next(nextUniverse)
		currentUniverse, nextUniverse = nextUniverse, currentUniverse

		update := frameTimer.End(start)
		if update {
			status = frameTimer.Status()
		}

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	js.Global().Set("startGame", js.FuncOf(InitGame))
	js.Global().Call("requestAnimationFrame", renderFrame)
	select {}
}
