package asn1per


const (
	ALIGNED		= iota
	UNALIGNED
)


func appendBit(b byte, bit byte) byte {
	b <<= 1
	b |= bit
	return b
}