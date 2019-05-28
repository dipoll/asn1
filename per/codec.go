package per

import (
	"errors"
	"math/big"
)

// BitEncoder implements basic operations needed to encode
// message in PER format
type BitEncoder struct {
	bits uint
	buf  *big.Int
}

// NewBitEncoder returns initialized new encoder
func NewBitEncoder() *BitEncoder {
	return &BitEncoder{buf: big.NewInt(0)}
}

// AppendBit appends bits to the left
func (e *BitEncoder) AppendBit(b uint) int {
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
func (e *BitEncoder) BitLen() uint {
	return e.bits
}

// Len returns current length in bytes to be written
func (e *BitEncoder) Len() uint {
	nbytes, nbits := e.FullLen()
	if nbits > 0 {
		nbytes++
	}
	return nbytes
}

// FullLen returns number of fully occupied bytes and number of occupied
// bits in the latest(right) byte
func (e *BitEncoder) FullLen() (uint, uint) {
	bytes := e.bits / 8
	trail := e.bits % 8
	return bytes, trail
}

// Bytes returns encoded bytes
func (e *BitEncoder) Bytes() []byte {
	return e.buf.Bytes()
}

// Align aligns bits to bytes, so the next
// encoded type will be padded to the next byte
func (e *BitEncoder) Align() {
	e.bits = e.Len() * 8
	e.buf = e.buf.Lsh(e.buf, 8)
}

// AppendInteger appends integer of defined bit length
// to the number
func (e *BitEncoder) AppendInt(num *big.Int, nBits uint) int {
	_, b := e.FullLen()
	shift := int(8-b) - int((nBits % 8))
	switch {
	case shift > 0:
		num = num.Lsh(num, uint(shift))
	case shift < 0:
		num = num.Lsh(num, 8)
		num = num.Rsh(num, uint(-1*shift))
	}
	pShift := (nBits - (8 - b) + 7) / 8
	if pShift > 0 {
		e.buf = e.buf.Lsh(e.buf, pShift*8)
	}

	e.buf = e.buf.Add(e.buf, num)
	e.bits += nBits
	return int(nBits)
}

// AppendBytes appends pure bytes to the end of buffer
func (e *BitEncoder) AppendBytes(b []byte) int {
	return e.AppendInt(big.NewInt(0).SetBytes(b), uint(len(b)*8))
}

// Reset clears encoder's buffer and bit counter. Should
// be called after each message has been encoded/decoded
func (e *BitEncoder) Reset() {
	e.buf = big.NewInt(0)
	e.bits = 0
}

// EncodeLength get encoded tag and length in chanks
// if length is greater than the biggest chunk remainin
// part of the length is returned
func EncodeLength(length int) (det []byte, remain int) {
	remain = 0
	switch {
	case length < 128:
		det = []byte{byte(length)}
	case length < 16384:
		num := uint16(length&0xFF | int(byte(0x80)|byte(length>>8))<<8)
		det = []byte{byte(num >> 8), byte((num << 8) >> 8)}
	case length < 32768:
		num := uint64(0xC1<<61) | uint64(length)
		det = []byte{byte(num >> 56),
			byte(num >> 48),
			byte(num >> 40),
			byte(num >> 32),
			byte(num >> 16),
			byte(num >> 8), byte(num)}
	}
	return
}

// ReadBit reads one bit at position pos
// returns 0 or 1 and no error in case of
// successful reading. In case of error returns
// -1 and error
func ReadBit(pos int, s []byte) (int, error) {
	b, _, err := ReadBits(pos, 1, s)
	if err == nil && b.Int64() > -1 {
		return int(b.Int64()), nil
	}
	return -1, err
}

// ReadBits reads big.Int number from the byte slice
// if there is no enough bytes available
// remaining bits and error will be returned as well as
// partilally read value
func ReadBits(pos int, length int, s []byte) (*big.Int, int, error) {
	startByte := pos / 8
	startBit := pos % 8
	out := big.NewInt(0)
	for i := startByte; i < len(s); i++ {
		if length == 0 {
			break
		}
		data := s[i]
		if startBit > 0 {
			data = data << startBit
			data = data >> startBit
		}
		chank := 8 - startBit
		if chank >= length {
			chank = length
		}
		out = out.Lsh(out, uint(chank))
		out = out.Add(out, big.NewInt(int64(data>>(8-(startBit+chank)))))
		length -= chank
		startBit = 0
	}
	if length > 0 {
		return out, length, errors.New("PER: Partial read, not enough bytes")
	}
	return out, length, nil
}
