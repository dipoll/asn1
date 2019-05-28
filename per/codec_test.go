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
	length int
	det    []byte
	remain int
}{
	{140, []byte{0x80, 0x8c}, 0},
	{67, []byte{0x43}, 0},
	{32000 , []byte{0x30,0x82,0x7D,0x04,0x80,0x82,0x7D,0x00}, 0}}

func TestEncodeLength(t *testing.T) {
	for n, v := range lenDetT {
		det, remain := EncodeLength(v.length)
		if !equal(det, v.det) || remain != v.remain {
			t.Errorf("TestEncodeLength [%d]: Expect Det: %08b Rem: %d, Got Det: %08b Rem: %d\n",
				n, v.det, v.remain, det, remain)
		}
	}
}
