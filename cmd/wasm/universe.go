package main

import (
	"math/rand"
	"unsafe"
)

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

	return rules
}
func initHighLife() []uint8 {
	// 35 = 0b100011
	rules := make([]uint8, 36)

	// 3 neighbours ALIVE should survive
	rules[3 << 2 | 1] = 1
	// 2 neighbours ALIVE should survive
	rules[2 << 2 | 1] = 1

	// 3 neighbours DEAD should be born
	rules[3 << 2 | 0] = 1
	// 6 neighbours DEAD should be born
	rules[6 << 2 | 0] = 1
	return rules
}
func initSeeds() []uint8 {
	// 35 = 0b100011
	rules := make([]uint8, 36)

	// 2 neighbours DEAD should be born
	rules[2 << 2 | 0] = 1
	return rules
}
func initBriansBrain() []uint8 {
	// 35 = 0b100011
	rules := make([]uint8, 36)

	// 2 neighbours DEAD should be born
	rules[2 << 2 | 0] = 1

	// all ALIVE should be dying
	rules[1] = 2
	rules[1 << 2 | 1] = 2
	rules[2 << 2 | 1] = 2
	rules[3 << 2 | 1] = 2
	rules[4 << 2 | 1] = 2
	rules[5 << 2 | 1] = 2
	rules[6 << 2 | 1] = 2
	rules[7 << 2 | 1] = 2
	rules[8 << 2 | 1] = 2

	// all DYING should die
	rules[2] = 0
	rules[1 << 2 | 2] = 0
	rules[2 << 2 | 2] = 0
	rules[3 << 2 | 2] = 0
	rules[4 << 2 | 2] = 0
	rules[5 << 2 | 2] = 0
	rules[6 << 2 | 2] = 0
	rules[7 << 2 | 2] = 0
	rules[8 << 2 | 2] = 0
	return rules
}

func initStarWars() []uint8 {
	rules := make([]uint8, 36)
	// B2/S345/C4
	rules[2 << 2 | 0] = 1

	rules[1] = 2
	rules[1 << 2 | 1] = 2
	rules[2 << 2 | 1] = 2
	
	rules[3 << 2 | 1] = 1
	rules[4 << 2 | 1] = 1
	rules[5 << 2 | 1] = 1

	rules[6 << 2 | 1] = 2
	rules[7 << 2 | 1] = 2
	rules[8 << 2 | 1] = 2

	rules[2] = 3
	rules[1 << 2 | 2] = 3
	rules[2 << 2 | 2] = 3
	rules[3 << 2 | 2] = 3
	rules[4 << 2 | 2] = 3
	rules[5 << 2 | 2] = 3
	rules[6 << 2 | 2] = 3
	rules[7 << 2 | 2] = 3
	rules[8 << 2 | 2] = 3

	rules[3] = 0
	rules[1 << 2 | 3] = 0
	rules[2 << 2 | 3] = 0
	rules[3 << 2 | 3] = 0
	rules[4 << 2 | 3] = 0
	rules[5 << 2 | 3] = 0
	rules[6 << 2 | 3] = 0
	rules[7 << 2 | 3] = 0
	rules[8 << 2 | 3] = 0
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
func isAlive(cell uint8) uint8 {
	r := cell & 1
	l := cell >> 1
	return r &^ l
}
func (u *Universe) Next(target *Universe) (*Universe, *Universe) {
	u.syncEdges()
	srcPtr := uintptr(unsafe.Pointer(&u.cells[0]))
	dstPtr := uintptr(unsafe.Pointer(&target.cells[0]))
    w := uintptr(u.width + 2)

	var offset, pA, pM, pB, tp uintptr
	var vA, vM, vB, self uint8
	var leftCol, midCol, rightCol, total uint8
    for y := 1; y <= u.height; y++ {

    	offset = uintptr(y) * w
        pA = srcPtr + offset - w
        pM = srcPtr + offset
        pB = srcPtr + offset + w
        
        tp = dstPtr + offset + 1

		vA = *(*uint8)(unsafe.Pointer(pA))
		vM = *(*uint8)(unsafe.Pointer(pM))
		vB = *(*uint8)(unsafe.Pointer(pB))
        leftCol = (vA & 1) &^ (vA >> 1) +
        	(vM & 1) &^ (vM >> 1) +
        	(vB & 1) &^ (vB >> 1)

		vA = *(*uint8)(unsafe.Pointer(pA + 1))
		vM = *(*uint8)(unsafe.Pointer(pM + 1))
		vB = *(*uint8)(unsafe.Pointer(pB + 1))
        midCol = (vA & 1) &^ (vA >> 1) +
        	(vM & 1) &^ (vM >> 1) +
        	(vB & 1) &^ (vB >> 1)

        pA += 2
        pM += 2
        pB += 2
        
        for x := 1; x <= u.width; x++ {
        	// self is previous vM
        	self = vM
        	vA = *(*uint8)(unsafe.Pointer(pA))
        	vM = *(*uint8)(unsafe.Pointer(pM))
        	vB = *(*uint8)(unsafe.Pointer(pB))
			rightCol = (vA & 1) &^ (vA >> 1) +
				(vM & 1) &^ (vM >> 1) +
				(vB & 1) &^ (vB >> 1)
			total = leftCol + midCol + rightCol

			*(*uint8)(unsafe.Pointer(tp)) =
				rulesTable[(total - (self & 1) &^ (self >> 1)) << 2 | self]

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
