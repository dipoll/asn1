package asn1per

import (
	"fmt"
	"math/bits"
)

// parseBool parser single boolean value for both
// aligned and unaligned per from single byte
func parseBool(offset uint8, bits byte) (bool, error) {
	pos := 1 << (7 - offset)
	if (bits & byte(pos)) > 0 {
		return true, nil
	}

	return false, nil
}

// appendBool adds boolean value to the single bit
func appendBool(bits *byte, offset uint8, v bool) uint8 {
	bit := byte(0)
	if v {
		bit = 1
	}
	*bits |= (bit << (7 - offset))
	return offset + 1
}

// Encoder represents
type Encoder struct {
	nbits int64 // Track current number of bits in encoded bytes sequence
	buf   []byte
}

func (e *Encoder) addUnsignedNumber(num uint64) error {
	fmt.Printf("Bits needed to write: %d = %d\n", num, bits.Len64(num))
	return nil
}
