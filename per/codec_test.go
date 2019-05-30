package per

import "testing"

// equal is helper function to compare byte
// slices
func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestEncoderInt(t *testing.T) {

}

var readBitsT = []struct {
	b   []byte
	pos int
	l   int
	v   int64
	err bool
}{
	{[]byte{128}, 0, 2, 2, false},
	{[]byte{130, 128}, 6, 4, 10, false},
	{[]byte{84, 255, 170}, 3, 14, 10751, false},
	{[]byte{84, 255, 170}, 15, 9, 426, false},
	{[]byte{84, 255, 170}, 15, 20, 426, true},
	{[]byte{}, 15, 20, 0, true}}

func TestReadBits(t *testing.T) {
	for n, v := range readBitsT {
		b, l, e := ReadBits(v.pos, v.l, v.b)
		if e != nil && !v.err {
			t.Errorf("TestReadBits [%d]: Got unexpected error: %v\n", n, e)
		} else if e == nil && v.err {
			t.Errorf("TestReadBits [%d]: Expect error but got nil\n", n)
		}
		if b.Int64() != v.v {
			if v.l != 0 && ((v.pos+v.l)-(len(v.b)*8)) != l {
				t.Errorf("TestReadBits [%d]: Incorrect remaining bits number:%d\n", n, l)
			}
			t.Errorf("TestReadBits [%d]: Expect: %d Got: %d\n", n, v.v, b.Int64())
			t.Errorf("TestReadBits [%d]: Dump\n Data-->%08b\n Got--->%08b\n", n, v.b, b.Int64())
		}
	}
}

var readBitT = []struct {
	b   []byte
	pos int
	v   int
	err bool
}{
	{[]byte{128}, 0, 1, false},
	{[]byte{127}, 0, 0, false},
	{[]byte{127}, 6, 1, false},
	{[]byte{127}, 7, 1, false},
	{[]byte{127}, 9, -1, true}}

func TestReadBit(t *testing.T) {
	for n, v := range readBitT {
		b, e := ReadBit(v.pos, v.b)
		if e != nil && !v.err {
			t.Errorf("TestReadBit [%d]: Got unexpected error: %v\n", n, e)
		} else if e == nil && v.err {
			t.Errorf("TestReadBit [%d]: Expect error but got nil\n", n)
		}
		if b != v.v {
			t.Errorf("TestReadBit [%d]: Expect: %d Got: %d\n", n, v.v, b)
			t.Errorf("TestReadBit [%d]: Dump\n Data-->%08b\n Got--->%08b\n", n, v.b, b)
		}
	}
}

var lenDetT = []struct {
	length   int
	det      []byte
	consumed int
}{
	{140, []byte{0x80, 0x8c}, 140},
	{67, []byte{0x43}, 67},
	{32000, []byte{0xC1}, 16384}}

func TestLengthDet(t *testing.T) {
	for n, v := range lenDetT {
		det, consumed := LengthDet(v.length)
		if !equal(det, v.det) || consumed != v.consumed {
			t.Errorf("TestEncodeLength [%d]: Expect Det: %08b Rem: %d, Got Det: %08b Rem: %d\n",
				n, v.det, v.consumed, det, consumed)
		}
	}
}

var readLenDetT = []struct{
	data	[]byte
	bRead	int
	length	int
	err		bool
}{
	{[]byte{0x80,0xFE,0x61,0x61}, 2*8, 254, false},
	{[]byte{0x81, 0x00, 0x61}, 2*8, 256, false},
	{[]byte{0x78,0x61}, 1*8, 120, false},
	{[]byte{0x80,0x80,0x61}, 2*8, 128, false},
	{[]byte{0xC4,0x61}, 1*8, 65536, false},
	{[]byte{0xC3,0x61}, 1*8, 49152, false},
	{[]byte{0x80}, 0, 0, true}}

func TestReadLenDet(t *testing.T){
	for n, v := range readLenDetT {
		rB, cS, err := ReadLenDet(0, v.data)
		if err == nil && v.err  {
			t.Errorf("TestReadLenDet [%d]: Should return error!", n)
			continue
		} else if err != nil && v.err {
			continue
		}
		if rB != v.bRead || cS != v.length {
			t.Errorf("TestReadLenDet [%d]: Decode: want: %d, got: %d", n, v.length, cS)
		}
	}
}