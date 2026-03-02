package main

import (
	"fmt"
)

type Color struct {
	R, G, B, A uint8
}

func (c Color) String() string {
	return fmt.Sprintf("rgba(%d, %d, %d, %f)",
		c.R, c.G, c.B, float64(c.A)/255.0)
}
func (c Color) Solid() Color {
	return Color{
		R: c.R,
		G: c.G,
		B: c.B,
		A: 0xFF,
	}
}
func (c *Color) Packed() uint32 {
	return uint32(c.R) |
		uint32(c.G)<<8 |
		uint32(c.B)<<16 |
		uint32(c.A)<<24
}
