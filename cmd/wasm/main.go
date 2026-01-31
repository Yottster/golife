package main

import (
	"fmt"
	"syscall/js"
//	"time"
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

type FPSCounter struct {
	lastTimestamp float64
	frameTimes []float64
	limit uint8
	writeIndex uint8
	frameCount uint8
}

func NewFPSCounter(limit uint8) *FPSCounter {
	return &FPSCounter {
		limit: limit,
		frameTimes: make([]float64, limit),
		writeIndex: 0,
		frameCount: 0,
	}
}

func (f *FPSCounter) Add(dt float64) {
	f.frameTimes[f.writeIndex] = dt
	f.writeIndex++
	f.writeIndex %= f.limit
	f.frameCount++
}

func (f *FPSCounter) Avg() float64 {
	var sum float64 = 0
	for _, value := range f.frameTimes {
		sum += value
	}
	return sum / float64(f.limit)
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

	var lastTimestamp float64
	var fps uint8

	fpsCounter := NewFPSCounter(30)
	performance := js.Global().Get("performance")
	ctx.Set("font", "40px monospace")
	var renderFrame js.Func
	var status string
	
	renderFrame = js.FuncOf(func (this js.Value, args []js.Value) interface {} {
		start := performance.Call("now").Float()
		Paint(ctx, currentUniverse, cellSize, status)
		
		currentUniverse.Next(nextUniverse)
		currentUniverse, nextUniverse = nextUniverse, currentUniverse

		end := performance.Call("now").Float()

		fpsCounter.Add(end - start)
		if end - lastTimestamp >= 1000 {
			fps = fpsCounter.frameCount
			fpsCounter.frameCount = 0
			lastTimestamp = end
			status = fmt.Sprintf(
				"FPS: %d | CPU: %.2f ms", fps, fpsCounter.Avg())
		}

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	js.Global().Set("startGame", js.FuncOf(InitGame))
	js.Global().Call("requestAnimationFrame", renderFrame)
	select {}
}
