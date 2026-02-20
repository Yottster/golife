package main

import (
	"math/rand"
	"unsafe"
)

var rulesTable []uint8

func init() {
	rulesTable = initGameOfLifeRules()
}
type Universe struct {
    width int
    height int
    cells []uint8
}
func NewUniverse(width int, height int) *Universe {
    var length int = (width + 2) * (height + 2)
    var universe Universe
    universe.width = width
    universe.height = height
    universe.cells = make([]uint8, length)

    return &universe
}
func (u *Universe) seed() {
    length := (u.width + 2) * (u.height + 2)
    for i := 0; i < length; i++ {
        u.cells[i] = uint8(rand.Intn(2))
    }
}
func initGameOfLifeRules() []uint8 {
	// 35 = 0b100011
	rules := make([]uint8, 36)

	// 3 neighbours ALIVE should survive
	rules[3 << 2 | 1] = 1
	// 2 neighbours ALIVE should survive
	rules[2 << 2 | 1] = 1
	// 3 neighbours DEAD should be born
	rules[3 << 2 | 0] = 1

	rules[6 << 2 | 0] = 1
	return rules
}

func (u *Universe) syncEdges() {
	w := u.width
	h := u.height
	rowWidth := w + 2

	// top-left
	u.cells[0] = u.cells[h*rowWidth + w]
	// top-right
	u.cells[rowWidth-1] = u.cells[h*rowWidth + 1]
	// bottom-left
	u.cells[(h+1) * rowWidth] = u.cells[rowWidth + w]
	// bottom-right
	u.cells[(h+1) * rowWidth + rowWidth - 1] = u.cells[rowWidth + 1]

	// copy top + bottom
	for x := 1; x <= w; x++ {
		u.cells[x] = u.cells[h*rowWidth + x]
		u.cells[(h+1)*rowWidth + x] = u.cells[x + rowWidth]
	}

	// copy left + right
	for y := 1; y <= h; y++ {
		u.cells[y*rowWidth] = u.cells[y*rowWidth + w]
		u.cells[y*rowWidth + w + 1] = u.cells[y*rowWidth + 1]
	}
}

func ruleCheck(neighbours uint8, cell uint8) uint8 {
	return rulesTable[neighbours << 2 | cell]
}

func (u *Universe) Next(target *Universe) (*Universe, *Universe) {
	u.syncEdges()
	srcPtr := uintptr(unsafe.Pointer(&u.cells[0]))
	dstPtr := uintptr(unsafe.Pointer(&target.cells[0]))
    w := uintptr(u.width + 2)

    for y := 1; y <= u.height; y++ {

    	offset := uintptr(y) * w
        pA := srcPtr + offset - w
        pM := srcPtr + offset
        pB := srcPtr + offset + w
        
        tp := dstPtr + offset + 1

        leftCol := uint8(
        	*(*uint8)(unsafe.Pointer(pA)) +
        	*(*uint8)(unsafe.Pointer(pM)) +
        	*(*uint8)(unsafe.Pointer(pB)))

        midCol := uint8(
        	*(*uint8)(unsafe.Pointer(pA + 1)) +
        	*(*uint8)(unsafe.Pointer(pM + 1)) +
        	*(*uint8)(unsafe.Pointer(pB + 1)))

        pA += 2
        pM += 2
        pB += 2
        
        for x := 1; x <= u.width; x++ {
			rightCol := uint8(
				*(*uint8)(unsafe.Pointer(pA)) +
				*(*uint8)(unsafe.Pointer(pM)) +
				*(*uint8)(unsafe.Pointer(pB)))
			total := leftCol + midCol + rightCol
			self := *(*uint8)(unsafe.Pointer(pM - 1))

			*(*uint8)(unsafe.Pointer(tp)) = ruleCheck(
				total - self, self)

			leftCol = midCol
			midCol = rightCol

			pA++
			pM++
			pB++
            tp++
        }
    }

    return target, u
}
