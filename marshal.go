package asn1per

import (
	"encoding/binary"
	"errors"
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
	offset      uint8 // Track current number of bits in encoded bytes sequence
	buf         []byte
	isAligned   bool
	isCanonical bool
}

// appendUint added unsigned integer 64 with number
// bits in it
func (e *Coder) appendUint64(v uint64, n uint8) error {
	if n > 64 {
		return errors.New("Number of bits greater than 64")
	}

	newNum := uint64(0)
	newNum = v << uint64(64-n)

	if  0 < e.offset && e.offset < 8 {
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

func (e *Coder) appendLenDeterminant(length uint64) (encoded uint64, nlength uint64) {
	switch {
	case length < 128:
		encoded = length
	case length < 16384:
		encoded = length&0xFF | uint64(byte(0x80)|byte(length>>8))<<8
	case length < 32768:
		encoded = 0xC1
		nlength = 16384
	case length < 49152:
		encoded = 0xC2
		nlength = 32768
	case length < 65536:
		encoded = 0xC3
		nlength = 49152
	default:
		encoded = 0xC4
		nlength = 65536

	}
	return
}

func (e *Coder) appendUint64Bytes(value uint64) int {
	l := binary.Size(value)
	e.appendUint64(value, uint8(l*8))
	return int(l * 8)
}

func (e *Coder) appendConstrainedUint64(value, min, max int64) int {
	rng := (max - min + 1)
	value -= min
	l := bits.Len64(uint64(rng))
	if rng > 255 && e.isAligned {
		if e.offset != 0 {
			e.buf = append(e.buf, byte(0))
			e.offset = 0
		}
	}
	switch {
	case rng <= 255:
		e.appendUint64(uint64(value), uint8(l))
		return int(rng)
	case rng == 256:
		e.appendUint64(uint64(value), 8)
		return 8
	case rng <= 65536:
		e.appendUint64(uint64(value), 16)
		return 16
	default:
		e.appendUint64(uint64(value), uint8(l))
	}
	return l
}

// BitLen returns encoded length in bits
func (e Coder) BitLen() int {
	return (len(e.buf)-1) * 8 + int(e.offset)
}
// addBool adds boolean value to the single bit
func (e *Coder) addBool(v bool) int {
	if v {
		e.appendBit(1)

	} else {
		e.appendBit(0)
	}
	return 1
}

func (e *Coder) appendBit(b byte) int {
	if e.offset > 7 || len(e.buf) < 1 {
		e.buf = append(e.buf, byte(0))
		e.offset = 0
	}
	e.buf[len(e.buf)-1] |= (b << (7 - e.offset))
	e.offset++
	return 1
}
