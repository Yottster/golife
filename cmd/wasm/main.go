package main

import (
	"syscall/js"
	"time"
	"unsafe"
)

func Render(u *Universe, dst []uint32, colors []uint32) {
	// potential BCE hint
	if len(colors) < 4 {
		return
	}

	w := u.width
	h := u.height
	_ = dst[:w*h]
	_ = u.cells[:(h+1)*(w+2)+w]

	// Build a full 256-entry lookup so uint8 indexing is always in bounds
	var lut [256]uint32
	lut[0] = colors[0]
	lut[1] = colors[1]
	lut[2] = colors[2]
	lut[3] = colors[3]

	rowWidth := w + 2

	for y := 0; y < h; y++ {
		rowStart := (y+1)*rowWidth + 1
		rowDst := dst[y*w : (y+1)*w : (y+1)*w]

		x := 0
		for ; x+7 < len(rowDst); x += 8 {
			cellIdx := rowStart + x
			chunk := *(*uint64)(unsafe.Pointer(&u.cells[cellIdx]))

			r := rowDst[x : x+8 : x+8]
			r[0] = lut[uint8(chunk)]
			r[1] = lut[uint8(chunk>>8)]
			r[2] = lut[uint8(chunk>>16)]
			r[3] = lut[uint8(chunk>>24)]
			r[4] = lut[uint8(chunk>>32)]
			r[5] = lut[uint8(chunk>>40)]
			r[6] = lut[uint8(chunk>>48)]
			r[7] = lut[uint8(chunk>>56)]
		}
		// rest
		for ; x < len(rowDst); x++ {
			rowDst[x] = lut[u.cells[rowStart+x]]
		}
	}
}
func themeChoice(choice int) []Color {
	choices := [][]Color{
		{
			// Amber
			{0, 0, 0, 0xFF},
			{0xFF, 0xEA, 0, 0xFF},
			{0xFF, 0x66, 0, 0xFF},
			{0x8B, 0, 0, 0xFF}},
		{
			// CRT Phosphor
			{0, 0, 0, 0xFF},
			{0, 0xFF, 0, 0xFF},
			{0, 0x88, 0, 0xFF},
			{0, 0x33, 0, 0xFF}},
		{
			// Deep Sea
			{0, 0, 0, 0xFF},
			{0, 0xFF, 0xFF, 0xFF},
			{0, 0x66, 0xFF, 0xFF},
			{0x4B, 0, 0x82, 0xFF}},
		{
			// star wars
			{0, 0, 0, 0xFF},
			{0, 0xFF, 0xFF, 0xFF},
			{0xFF, 0, 0x55, 0xFF},
			{0x33, 0, 0x33, 0xFF},
		},
	}
	return choices[choice%len(choices)]
}

// TODO: We should remove this global
var rulesTable [64]uint8

func main() {
	window := js.Global()
	document := window.Get("document")

	body := document.Get("body")
	bodyStyle := body.Get("style")

	theme := themeChoice(0)
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

	pixels := make([]byte, 4*width*height)
	ptr := unsafe.Pointer(&pixels[0])

	gameObject.Get("fn").Call("setMemoryView", uintptr(ptr))
	u32Buffer :=
		unsafe.Slice((*uint32)(ptr), width*height)

	nextUniverse := NewUniverse(width, height)
	currentUniverse := NewUniverse(width, height)
	currentUniverse.seed()

	timingsCollector := gameObject.Get("fn").Get("addGoTimings")

	var tick js.Func = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
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

		rulesets := []func() [64]uint8{
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
