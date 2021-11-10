package sdo

func hasBit(n uint8, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func setBit(n uint8, pos uint) uint8 {
	n |= (1 << pos)
	return n
}
