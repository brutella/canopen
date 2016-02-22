package canopen

type Index struct {
	B0 byte
	B1 byte
}

type ObjectIndex struct {
	Index    Index
	SubIndex byte
}

func NewObjectIndex(index uint16, subIndex uint8) ObjectIndex {
	return ObjectIndex{
		Index: Index{
			B0: byte(index & 0xFF),
			B1: byte(index >> 8),
		},
		SubIndex: subIndex,
	}
}
