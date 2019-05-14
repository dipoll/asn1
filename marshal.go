package asn1per

import (
	"encoding/binary"
	"errors"
	"fmt"
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

// printBytes
func printBytes(b []byte) {
	for n, v := range b {
		fmt.Printf("%08b ", v)
		if n%2 > 0 {
			fmt.Println("")
		}
	}
	fmt.Println("")
}

func NewEncoder() *Coder {
	return &Coder{
		buf: []byte{0}}
}

// Coder represents
type Coder struct {
	offset    uint8 // Track current number of bits in encoded bytes sequence
	buf       []byte
	isAligned bool
}

// addUint added unsigned integer 64 with number
// bits in it
func (e *Coder) addUint(v uint64, n uint8) error {
	if n > 64 {
		return errors.New("Number of bits greater than 64")
	}

	newNum := uint64(0)
	newNum = v<<64 - n

	if e.offset < 8 {
		newNum = newNum >> e.offset
	}

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, newNum)

	if e.offset < 8 {
		e.buf[len(e.buf)-1] |= buf[0]
		buf = buf[1:]
	}
	e.offset += n
	if e.offset <= 8 {
		return nil
	}

	nBytes := (e.offset / 8)
	if nBytes < 1 {
		return nil
	}
	e.offset = e.offset % 8
	if e.offset > 0 {
		nBytes++
	}

	e.buf = append(e.buf, buf[:nBytes-1]...)
	return nil
}

// addBool adds boolean value to the single bit
func (e *Coder) addBool(v bool) int {
	bit := byte(0)
	if v {
		bit = 1
	}
	fmt.Printf("BOOL: Offset %d, Adding Boolean %v\n", e.offset, v)
	if e.offset > 7 || len(e.buf) < 1 {
		e.buf = append(e.buf, byte(0))
		e.offset = 0
	}
	e.buf[len(e.buf)-1] |= (bit << (7 - e.offset))
	e.offset++
	return 1
}
