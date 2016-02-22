package sdo

// splitN splits b into a list of n sized bytes
func splitN(b []byte, n int) [][]byte {
	if len(b) < n {
		return [][]byte{b}
	}

	var bs [][]byte
	var buf []byte
	for i := 0; i < len(b); i++ {
		if len(buf) == n {
			bs = append(bs, buf)
			buf = []byte{}
		}

		buf = append(buf, b[i])
	}

	if len(buf) > 0 {
		bs = append(bs, buf)
	}

	return bs
}
