package main

import (
	"math/rand"
	"unsafe"
)

type Universe struct {
	width  int
	height int
	cells  []uint8
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
	cells := u.cells[:(u.width+2)*(u.height+2)]
	for i := range cells {
		cells[i] = uint8(rand.Intn(2))
	}
}
func (u *Universe) syncEdges() {
	w := u.width
	h := u.height
	rowWidth := w + 2
	lastRow := (h + 1) * rowWidth

	cells := u.cells[:lastRow+rowWidth]

	// corners
	cells[0] = cells[h*rowWidth+w]
	cells[w+1] = cells[h*rowWidth+1]
	cells[lastRow] = cells[rowWidth+w]
	cells[lastRow+w+1] = cells[rowWidth+1]

	// Top + bottom edges
	topRow := cells[1 : w+1]
	bottomRow := cells[lastRow+1 : lastRow+w+1]
	srcTop := cells[rowWidth+1 : rowWidth+w+1]
	srcBottom := cells[h*rowWidth+1 : h*rowWidth+w+1]
	copy(topRow, srcBottom)
	copy(bottomRow, srcTop)

	// Left + right edges
	for y := 1; y <= h; y++ {
		off := y * rowWidth
		cells[off] = cells[off+w]
		cells[off+w+1] = cells[off+1]
	}
}

func (u *Universe) Next(target *Universe) (*Universe, *Universe) {
	u.syncEdges()

	w := u.width + 2
	totalLen := (u.height + 2) * w

	src := unsafe.Slice(&u.cells[0], totalLen)
	dst := unsafe.Slice(&target.cells[0], totalLen)

	for y := 1; y <= u.height; y++ {
		rows := src[(y-1)*w : (y+2)*w]
		rowDst := dst[y*w+1 : (y+1)*w-1 : (y+1)*w-1]
		n := len(rowDst)

		aLeft := rows[0 : n+2 : n+2]
		mLeft := rows[w : w+n+2 : w+n+2]
		bLeft := rows[2*w : 2*w+n+2 : 2*w+n+2]

		aRight := aLeft[2 : n+2 : n+2]
		mRight := mLeft[2 : n+2 : n+2]
		bRight := bLeft[2 : n+2 : n+2]

		vA := aLeft[0]
		vM := mLeft[0]
		vB := bLeft[0]
		leftCol := (vA&1)&^(vA>>1) +
			(vM&1)&^(vM>>1) +
			(vB&1)&^(vB>>1)

		vA = aLeft[1]
		vM = mLeft[1]
		vB = bLeft[1]
		midCol := (vA&1)&^(vA>>1) +
			(vM&1)&^(vM>>1) +
			(vB&1)&^(vB>>1)

		for x := range n {
			self := vM
			vA = aRight[x]
			vM = mRight[x]
			vB = bRight[x]
			rightCol := (vA&1)&^(vA>>1) +
				(vM&1)&^(vM>>1) +
				(vB&1)&^(vB>>1)
			total := leftCol + midCol + rightCol

			idx := (total-(self&1)&^(self>>1))<<2 | self
			rowDst[x] = rulesTable[idx&0x3F]

			leftCol = midCol
			midCol = rightCol
		}
	}

	return target, u
}
