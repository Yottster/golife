package main

import (
	"fmt"
	"syscall/js"
)

type FrameTimer struct {
	performance js.Value
	limit uint8
	frameTimes []float64
	writeIndex uint8
	count uint8
	lastTimestamp float64
	frameCount uint8
	recordedFPS uint8
}

func NewFrameTimer(performance js.Value, limit uint8) *FrameTimer {
	return &FrameTimer {
		performance: performance,
		limit: limit,
		frameTimes: make([]float64, limit),
		writeIndex: 0,
		count: 0,
		frameCount: 0,
		recordedFPS: 0,
	}
}

func (f *FrameTimer) Start() float64 {
	return f.performance.Call("now").Float()
}

func (f *FrameTimer) End(
	start float64) bool {
	currentTimestamp := 
		f.performance.Call("now").Float()
	f.frameTimes[f.writeIndex] = 
		currentTimestamp - start
	f.writeIndex++
	f.writeIndex %= f.limit
	f.frameCount++

	if (f.count < f.limit) {
		f.count++
	}
	

	if currentTimestamp - f.lastTimestamp >= 1000 {
		f.recordedFPS = f.frameCount
		f.frameCount = 0
		f.lastTimestamp = currentTimestamp
		return true
	}
	return false
}

func (f *FrameTimer) Status() string {
	return fmt.Sprintf(
			"FPS: %d | CPU: %.2f ms", 
			f.recordedFPS, 
			f.Avg())
}

func (f *FrameTimer) Avg() float64 {
	var sum float64 = 0
	for _, value := range f.frameTimes {
		sum += value
	}
	return sum / float64(f.count)
}
