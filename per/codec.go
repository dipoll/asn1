package per

import (
	"errors"
	"fmt"
	"math/big"
)

const (
	chunk16K = 1024 * 16
	chunk32K = chunk16K * 2
	chunk48K = chunk16K * 3
	chunk64K = chunk16K * 4
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
func (e *BitEncoder) AppendInt(num *big.Int, nBits int) int {
	_, b := e.FullLen()
	shift := int(8-b) - (nBits % 8)
	switch {
	case shift > 0:
		num = num.Lsh(num, uint(shift))
	case shift < 0:
		num = num.Lsh(num, 8)
		num = num.Rsh(num, uint(-1*shift))
	}
	pShift := (nBits - int(8-b) + 7) / 8
	if pShift > 0 {
		e.buf = e.buf.Lsh(e.buf, uint(pShift*8))
	}

	e.buf = e.buf.Add(e.buf, num)
	e.bits += uint(nBits)
	return nBits
}

// AppendBytes appends pure bytes to the end of buffer
func (e *BitEncoder) AppendBytes(b []byte) int {
	return e.AppendInt(big.NewInt(0).SetBytes(b), len(b)*8)
}

// Reset clears encoder's buffer and bit counter. Should
// be called after each message has been encoded/decoded
func (e *BitEncoder) Reset() {
	e.buf = big.NewInt(0)
	e.bits = 0
}

// AppendConstInt appends constrained integer to the byte buffer.
func (e *BitEncoder) AppendConstInt(value *big.Int, min, max int, align bool) int {
	rng := max - min + 1
	value = value.Add(value, big.NewInt(int64(min)))

	if rng > 255 {
		e.Align()
	}

	switch {
	case rng < 256:
		return e.AppendInt(value, value.BitLen())

	case rng == 256:
		return e.AppendInt(value, 8)

	case rng <= 65536:
		return e.AppendInt(value, 16)

	default:
		return e.AppendInt(value, value.BitLen())
	}
}

// AppendUnconstInt appends unconstrained signed integer
// to the byte buffer
func (e *BitEncoder) AppendUnsconstInt(v *big.Int) int {
	return 0
}

// AppendLenDet appends length determinant to the bytes
// to the internal buffer
func (e *BitEncoder) AppendLenDet(v *big.Int, length int) (nBits int, err error) {
	return
}

// LengthDet returns determinant encdoed as slice of
// bytes and consumed length by chunk. In the case if
// length could not be encoded into one chunk of
// data this function should be invoked until
// consumed == length
func LengthDet(length int) (det []byte, consumed int) {
	switch {
	case length < 128:
		det = []byte{byte(length)}
		consumed = length

	case length < chunk16K:
		num := uint16(length&0xFF | int(byte(0x80)|byte(length>>8))<<8)
		det = []byte{byte(num >> 8), byte((num << 8) >> 8)}
		consumed = length

	case length < chunk32K:
		det = []byte{0xC1}
		consumed = chunk16K

	case length < chunk48K:
		det = []byte{0xC2}
		consumed = chunk16K

	case length < chunk64K:
		det = []byte{0xC3}
		consumed = chunk48K

	default:
		det = []byte{0xC4}
		consumed = chunk64K
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

// ToNegative converts big.Int to negative
// representation for PER encoding
func ToNegative(v *big.Int) *big.Int {
	v = v.Xor(v, v)
	v = v.Add(v, big.NewInt(1))
	l := (v.BitLen() + 7) / 8
	fmt.Printf("Int before sum: %08b with Length: %d\n", v.Bytes(), l)
	sh := big.NewInt(1).Lsh(big.NewInt(1),uint(l*8))
	fmt.Printf("Shifted: %08b\n", sh.Bytes())
	v = v.Add(v, sh)
	fmt.Printf("Int after sum: %08b with Length: %d\n", v.Bytes(), (v.BitLen() + 7) / 8)
	z := v.And(v, big.NewInt(int64(1<<((8*l)-1))))
	if z.Cmp(big.NewInt(0)) == 0 {
		fmt.Printf("Before add %08b\n", z.Bytes())
		v = v.Or(v, big.NewInt(0xff<<(8*l)))
		l++
	}
	fmt.Printf("%08b\n", v)
	return v
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
		chunk := 8 - startBit
		if chunk >= length {
			chunk = length
		}
		out = out.Lsh(out, uint(chunk))
		out = out.Add(out, big.NewInt(int64(data>>(8-(startBit+chunk)))))
		length -= chunk
		startBit = 0
	}
	if length > 0 {
		return out, length, errors.New("per: codec: partial read, not enough bytes")
	}
	return out, length, nil
}

// ReadLenDet reads length determinant and returns
// read bits consumed by length determinant and size of data chunk to be read in bytes
func ReadLenDet(pos int, s []byte) (readBits, chunkSize int, err error) {
	det, readBits, err := ReadBits(pos, 16, s)
	if err != nil {
		return
	}

	bts := det.Bytes()
	readBits = 8

	switch {
	case (bts[0] & 0x80) == 0x00:
		chunkSize = int(bts[0])

	case (bts[0] & 0xC0) == 0x80:
		readBits = 16
		chunkSize = (int(bts[0]&0x7F) << 8) | int(bts[1])

	case bts[0] == 0xC1:
		chunkSize = chunk16K

	case bts[0] == 0xC2:
		chunkSize = chunk32K

	case bts[0] == 0xC3:
		chunkSize = chunk48K

	case bts[0] == 0xC4:
		chunkSize = chunk64K

	default:

		err = errors.New("per: ReadLenDet: can't read length determinant")

	}

	return
}
