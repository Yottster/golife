package main

const (
	dead   = 0
	alive  = 1
	dying  = 2
	fading = 3
)

func initGameOfLifeRules() [64]uint8 {
	// 35 = 0b100011
	var rules [64]uint8

	// 3 neighbours ALIVE should survive
	rules[3<<2|alive] = alive
	// 2 neighbours ALIVE should survive
	rules[2<<2|alive] = alive

	// 3 neighbours DEAD should be born
	rules[3<<2|dead] = alive

	return rules
}
func initHighLife() [64]uint8 {
	// 35 = 0b100011
	var rules [64]uint8

	// 3 neighbours ALIVE should survive
	rules[3<<2|alive] = alive
	// 2 neighbours ALIVE should survive
	rules[2<<2|alive] = alive

	// 3 neighbours DEAD should be born
	rules[3<<2|dead] = alive
	// 6 neighbours DEAD should be born
	rules[6<<2|dead] = alive
	return rules
}
func initSeeds() [64]uint8 {
	// 35 = 0b100011
	var rules [64]uint8

	// 2 neighbours DEAD should be born
	rules[2<<2|dead] = alive
	return rules
}
func initBriansBrain() [64]uint8 {
	// 35 = 0b100011
	var rules [64]uint8

	// 2 neighbours DEAD should be born
	rules[2<<2|dead] = alive

	// all ALIVE should be dying
	rules[alive] = dying
	rules[1<<2|alive] = dying
	rules[2<<2|alive] = dying
	rules[3<<2|alive] = dying
	rules[4<<2|alive] = dying
	rules[5<<2|alive] = dying
	rules[6<<2|alive] = dying
	rules[7<<2|alive] = dying
	rules[8<<2|alive] = dying

	// all DYING should die
	rules[dying] = dead
	rules[1<<2|dying] = dead
	rules[2<<2|dying] = dead
	rules[3<<2|dying] = dead
	rules[4<<2|dying] = dead
	rules[5<<2|dying] = dead
	rules[6<<2|dying] = dead
	rules[7<<2|dying] = dead
	rules[8<<2|dying] = dead
	return rules
}

func initStarWars() [64]uint8 {
	var rules [64]uint8
	// B2/S345/C4
	rules[2<<2|dead] = alive

	rules[alive] = dying
	rules[1<<2|alive] = dying
	rules[2<<2|alive] = dying

	rules[3<<2|alive] = alive
	rules[4<<2|alive] = alive
	rules[5<<2|alive] = alive

	rules[6<<2|alive] = dying
	rules[7<<2|alive] = dying
	rules[8<<2|alive] = dying

	rules[2<<2|dying] = fading
	rules[1<<2|dying] = fading
	rules[2<<2|dying] = fading
	rules[3<<2|dying] = fading
	rules[4<<2|dying] = fading
	rules[5<<2|dying] = fading
	rules[6<<2|dying] = fading
	rules[7<<2|dying] = fading
	rules[8<<2|dying] = fading

	rules[fading] = dead
	rules[1<<2|fading] = dead
	rules[2<<2|fading] = dead
	rules[3<<2|fading] = dead
	rules[4<<2|fading] = dead
	rules[5<<2|fading] = dead
	rules[6<<2|fading] = dead
	rules[7<<2|fading] = dead
	rules[8<<2|fading] = dead
	return rules
}
