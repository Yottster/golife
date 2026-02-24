package main

import (
	"syscall/js"
	"time"
	"unsafe"
)

func Render(u *Universe, dst []uint32, colors []uint32) {

	dstPtr := uintptr(unsafe.Pointer(&dst[0]))
	colorPtr := uintptr(unsafe.Pointer(&colors[0]))

	cellsPtr := uintptr(unsafe.Pointer(&u.cells[0]))
	rowWidth := uintptr(u.width + 2)
	
	step := unsafe.Sizeof(uint32(0))

	for y := 1; y <= u.height; y++ {
		rowOffset := uintptr(y) * rowWidth + 1
		cellPtr := cellsPtr + rowOffset

		for x := 1; x <= u.width; x++ {
			cell := *(*uint8)(unsafe.Pointer(cellPtr))
			c := *(*uint32)(unsafe.Pointer(
				colorPtr + uintptr(cell) * step))
			*(*uint32)(unsafe.Pointer(dstPtr)) = c
			cellPtr++
			dstPtr += step
		}
	}
}
func themeChoice(choice int) []Color {
	choices := [][]Color{
		[]Color{
			// Amber
			Color{ 0, 0, 0, 0xFF },
			Color{ 0xFF, 0xEA, 0, 0xFF },
			Color{ 0xFF, 0x66, 0, 0xFF },
			Color{ 0x8B, 0, 0, 0xFF }},
		[]Color{
			// CRT Phosphor
			Color{ 0, 0, 0, 0xFF },
			Color{ 0, 0xFF, 0, 0xFF },
			Color{ 0, 0x88, 0, 0xFF },
			Color{ 0, 0x33, 0, 0xFF }},
		[]Color{
			// Deep Sea
			Color{ 0, 0, 0, 0xFF },
			Color{ 0, 0xFF, 0xFF, 0xFF },
			Color{ 0, 0x66, 0xFF, 0xFF },
			Color{ 0x4B, 0, 0x82, 0xFF }},
		[]Color{
			// star wars
			Color{ 0, 0, 0, 0xFF},
			Color{ 0, 0xFF, 0xFF, 0xFF},
			Color{ 0xFF, 0, 0x55, 0xFF},
			Color{ 0x33, 0, 0x33, 0xFF},
		},
	}
	return choices[choice % len(choices)]
}

// TODO: We should remove this global
var rulesTable []uint8
func main() {
	window := js.Global()
	document := window.Get("document")

	body := document.Get("body")
	bodyStyle := body.Get("style")

	theme := themeChoice(3)
	bodyStyle.Set("background", theme[0].Solid().String())

    colors := make([]uint32, len(theme))
    for i, elem := range theme {
    	colors[i] = elem.Packed()
    }

	gameObject := window.Get("gol")
	dims := gameObject.Get("dims")

	width := dims.Get("width").Int()
	height := dims.Get("height").Int()

	println(width, height)
	
	pixels := make([]byte, 4 * width * height)
	ptr := unsafe.Pointer(&pixels[0])
	
	gameObject.Get("fn").Call("setMemoryView", uintptr(ptr))
	u32Buffer := 
		unsafe.Slice((*uint32)(ptr), width * height)

	nextUniverse := NewUniverse(width, height)
	currentUniverse := NewUniverse(width, height)
	currentUniverse.seed()

	timingsCollector := gameObject.Get("fn").Get("addGoTimings")

	var tick js.Func = js.FuncOf(func (this js.Value, args []js.Value) interface {} {
		t0 := time.Now()
		Render(currentUniverse, u32Buffer, colors)
	
		t1 := time.Now()
		renderTime := t1.Sub(t0).Microseconds()

		currentUniverse, nextUniverse =
			currentUniverse.Next(nextUniverse)
		nextTime := time.Since(t1).Microseconds()

		timingsCollector.Invoke(int(nextTime), int(renderTime))
		return nil
	})
	js.Global().Set("tick", tick)

	var setMode js.Func = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		mode := args[0].Int()

		rulesets := []func () []uint8 {
			initGameOfLifeRules,
			initHighLife,
			initSeeds,
			initBriansBrain,
			initStarWars,
		}

		if mode >= 0 && mode < len(rulesets) {
			rulesTable = rulesets[mode]()
		}
		return nil
	})
	gameObject.Get("fn").Set("setMode", setMode)

	select {}
}
