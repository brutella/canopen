package sdo

import "github.com/FabianPetersen/canopen"

func hasBit(n uint8, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func setBit(n uint8, pos uint) uint8 {
	n |= (1 << pos)
	return n
}

func getAbortCodeBytes(frame canopen.Frame) []uint8 {
	if len(frame.Data) >= 8 {
		return frame.Data[4:]
	}
	return []uint8{}
}
