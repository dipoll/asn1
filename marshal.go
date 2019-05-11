package asn1per

import (
	"encoding/binary"
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
func NewEncoder() *Coder {
	return &Coder{
		buf: []byte{0}}
}

// Coder represents
type Coder struct {
	offset    uint64 // Track current number of bits in encoded bytes sequence
	buf       []byte
	isAligned bool
}

// addUint64 appends uint64 number to the bytes
func (e *Coder) addUint64(num uint64, numBits uint64) error {
	tail := uint64(e.buf[len(e.buf)-1])
	tail <<= numBits
	tail |= num

	fmt.Printf("GG: %v\n", bits.Len64(tail))

	newOffset := (e.offset + numBits) % 8
	numBytes := (e.offset + numBits) / 8

	newTail := make([]byte, 4)
	fmt.Printf("NB: %v, OFF: %v, TLEN: %v, LEN: %v\n", numBytes, newOffset, len(newTail), len(e.buf))
	binary.PutUvarint(newTail, tail)
	e.buf = append(e.buf[:len(e.buf)-1], newTail[:3]...)
	fmt.Printf("FULL NEW LEN: %v\n", len(e.buf))
	e.offset = newOffset

	return nil
}

// addBool adds boolean value to the single bit
func (e *Coder) addBool(v bool) int {
	bit := byte(0)
	if v {
		bit = 1
	}
	if e.offset > 7 || len(e.buf) < 1 {
		e.buf = append(e.buf, byte(0))
	}
	fmt.Println(len(e.buf))
	e.buf[len(e.buf)-1] |= (bit << (7 - e.offset))
	e.offset++
	return 1
}
