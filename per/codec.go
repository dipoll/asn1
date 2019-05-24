package main

import (
	"fmt"
	"math/big"
)

type Encoder struct {
	bits uint
	buf  *big.Int
}

func NewEncoder() *Encoder {
	return &Encoder{buf: big.NewInt(0)}
}

// AppendBit appends bits to the left
func (e *Encoder) AppendBit(b uint) int {
	nbytes, nbits := e.FullLen()
	if nbits > 0 {
		nbytes++
	}
	e.buf.SetBit(e.buf, int(7-nbits), b)

	if nbits == 7 {
		e.buf = e.buf.Lsh(e.buf, 8)
	}
	e.bits++
	return 1
}

// BitLen returns current length in bits
func (e *Encoder) BitLen() uint {
	return e.bits
}

// Len returns current length in bytes to be written
func (e *Encoder) Len() uint {
	nbytes, nbits := e.FullLen()
	if nbits > 0 {
		nbytes++
	}
	return nbytes
}

// FullLen returns number of fully occupied bytes and number of occupied
// bits in the latest(right) byte
func (e *Encoder) FullLen() (uint, uint) {
	bytes := e.bits / 8
	trail := e.bits % 8
	return bytes, trail
}

// Bytes returns encoded bytes
func (e *Encoder) Bytes() []byte {
	return e.buf.Bytes()
}

// Align aligns bits to bytes, so the next
// encoded type will be padded to the next byte
func (e *Encoder) Align() {
	e.bits = e.Len() * 8
	e.buf = e.buf.Lsh(e.buf, 8)
}

// AppendInteger appends integer of defined bit length
// to the number
func (e *Encoder) AppendInteger(num big.Int, nBytes uint) int {
	return 1
}

// AppendBytes appends pure bytes to the end of buffer
func (e *Encoder) AppendBytes(b []bytes) int {
	return 1
}

func main() {
	encoder := NewEncoder()
	for i := 0; i < 15; i++ {
		encoder.AppendBit(1)
	}
	fmt.Println("Length before: ", encoder.bits)
	encoder.Align()
	fmt.Println("Length after: ", encoder.bits)
	encoder.AppendBit(0)
	encoder.AppendBit(1)
	fmt.Printf("Array: %08b\n", encoder.Bytes())
	fmt.Printf("Array: %d\n", encoder.bits)
}
