package main

import (
	"syscall/js"
	"strconv"
	"unsafe"
)

var privateUniverse *Universe

func InitGame(this js.Value, args []js.Value) interface{} {
	width := args[0].Int()
	height := args[1].Int()

	privateUniverse = NewUniverse(width, height)
	privateUniverse.seed()
	
	return nil
}

func Paint(
	target js.Value,
	writeCtx js.Value,
	pixels []byte) {

	js.CopyBytesToJS(target.Get("data"), pixels)
	writeCtx.Call("putImageData", target, 0, 0)
}
func TransferImage(
	srcCanvas js.Value,
	ctx js.Value,
	width int,
	height int) {

	ctx.Call("drawImage", srcCanvas, 0, 0, width, height)
}

func Render(u *Universe, dst []uint32, d uint32, l uint32) {
	
	for i, cell := range u.cells {
		if cell == 1 {
			dst[i] = l
		} else {
			dst[i] = d
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

	body := document.Get("body")
	bodyStyle := body.Get("style")

	deadColor := Color{0,0,0,200}
	liveColor := Color{0,0,0xFF,0xFF}

	deadPacked := deadColor.Packed()
	livePacked := liveColor.Packed()

	bodyStyle.Set("background", deadColor.Solid().String())
	bodyStyle.Set("margin", "0")
	bodyStyle.Set("padding", "0")
	
	canvas := document.Call("getElementById", "canvas")
	canvasStyle := canvas.Get("style")
	canvasStyle.Set("display", "block")
	canvasStyle.Set("overflow", "hidden")
	
	ctx := canvas.Call("getContext", "2d")

	location := window.Get("location")
	search := location.Get("search").String()
	params := window.Get("URLSearchParams").New(search)
	sizeParam := params.Call("get", "cellSize")
	var cellSize int = 3
	if !sizeParam.IsNull() {
		val, err := strconv.Atoi(sizeParam.String())
		if err == nil {
			cellSize = val
		}
	}
	
	innerWidth := window.Get("innerWidth").Int()
	innerHeight := window.Get("innerHeight").Int()

	canvas.Set("width", innerWidth)
	canvas.Set("height", innerHeight)

	width := innerWidth/cellSize
	height := innerHeight/cellSize

	hiddenCanvas := document.Call(
		"createElement", "canvas")
	hiddenCanvas.Set("width", width)
	hiddenCanvas.Set("height", height)
	writeCtx := hiddenCanvas.Call("getContext", "2d")
	
	pixels := make([]byte, 4 * width * height)
	ptr := unsafe.Pointer(&pixels[0])
	u32Buffer := 
		unsafe.Slice((*uint32)(ptr), width * height)
	
	imageData := writeCtx.Call(
		"createImageData", width, height)

	nextUniverse := NewUniverse(width, height)
	currentUniverse := NewUniverse(width, height)
	currentUniverse.seed()

	performance := window.Get("performance")
	frameTimer := NewFrameTimer(performance, 30)
	
	ctx.Set("font", "40px monospace")
	ctx.Set("imageSmoothingEnabled", false)
	var renderFrame js.Func
	var status string
	
	renderFrame = js.FuncOf(func (this js.Value, args []js.Value) interface {} {

		start := frameTimer.Start()

		Render(currentUniverse, u32Buffer, deadPacked, livePacked)

		Paint(imageData, writeCtx, pixels)
		TransferImage(hiddenCanvas, ctx, innerWidth, innerHeight)
		PrintStatus(ctx, status)

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
