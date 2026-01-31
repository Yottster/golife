package main

import (
	"math/rand"
)

type Universe struct {
    width int
    height int
    cells []uint8
}
func NewUniverse(width int, height int) *Universe {
    var length int = width * height
    var universe Universe
    universe.width = width
    universe.height = height
    universe.cells = make([]uint8, length)

    return &universe
}
func (u *Universe) seed() {
    length := u.width * u.height
    for i := 0; i < length; i++ {
        u.cells[i] = uint8(rand.Intn(2))
    }
}
func (u *Universe) index(x int, y int) int {
    return u.width * y + x
}
func (u *Universe) Alive(x int, y int) bool {
    x = (x + u.width) % u.width
    y = (y + u.height) % u.height
    return u.cells[u.index(x, y)] == 1
}

func (u *Universe) value(x int, y int) uint8 {
    x = (x + u.width) % u.width
    y = (y + u.height) % u.height
    return u.cells[u.index(x,y)]
}

func (u *Universe) Neighbours(x int, y int) uint8 {
    var sum uint8 = 0
    for xx := x - 1; xx < x + 2; xx++ {
        for yy := y - 1; yy < y + 2; yy++ {
            sum += u.value(xx, yy)
        }
    }
    return sum - u.value(x, y)
}

func rules(neighbours uint8, alive bool) uint8 {
    if neighbours == 3 {
        return 1
    }
    if neighbours < 2 || neighbours > 3 {
        return 0
    }
    if alive {
        return 1
    }
    return 0
}

func (u *Universe) Next(target *Universe) {

    i := 0
    for y := 0; y < u.height; y++ {
        for x := 0; x < u.width; x++ {
            neighbours := u.Neighbours(x, y)
            alive := u.Alive(x, y)
            target.cells[i] = rules(neighbours, alive)
            i++
        }
    }
}
