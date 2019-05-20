package asn1per

import (
	"testing"
)

var (
	testBool        byte   = 0xBF                     // PER and UPER in one byte {1 0 1 1 1 1 1 1}
	boolNumPER      []byte = []byte{0x80, 0x01, 0x02} // True and Unconstrained Number=2
	boolNumUPER     []byte = []byte{0x80, 0x81, 0x00}
	boolNumConstPER []byte = []byte{}
)

type tBoolNum struct {
	isA      bool
	Boolean1 bool
	Integer  uint64
	bytes    []byte
}

type tUnConstNum struct {
	isA     bool
	Integer []uint64
	Bytes   []byte
}

var tbBoolNum = []tBoolNum{
	tBoolNum{true, true, 2, []byte{0x80, 0x01, 0x02}},
	tBoolNum{false, true, 2, []byte{0x80, 0x81, 0x00}}}

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

func TestBooleanParsing(t *testing.T) {
	testBools := []bool{true, false, true, true,
		true, true, true, true}

	for i, v := range testBools {
		b, _ := parseBool(uint8(i), testBool)
		if b != v {
			t.Errorf("Bool Parse: POS: %d Expected:%v Got: %v\n", i, v, b)
		}
	}

	npos := uint8(0)
	out := byte(0)
	for _, v := range testBools {
		npos = appendBool(&out, npos, v)
	}
	if out != testBool {
		t.Errorf("Bool Encode Failed: Expect: %X, Got: %X", testBool, out)
	}
}

var tbConstrNumA = []struct {
	V    int64
	Min  int64
	Max  int64
	AL   int
	UL   int
	Ref  []byte
	RefU []byte
}{
	{-3, -5, 0, 8, 3, []byte{0x40}, []byte{0x40}},
	{127, 0, 255, 16, 11, []byte{0x40, 0x7F}, []byte{0x4F, 0xE0}},
	{256, 0, 256, 32, 20, []byte{0x40, 0x7F, 0x01, 0x00}, []byte{0x4F, 0xF0, 0x00}},
	{-72, -6900, 6546, 48, 34, []byte{0x40, 0x7F, 0x01, 0x00, 0x1A, 0xAC}, []byte{0x4F, 0xF0, 0x06, 0xAB, 0x00}}}

func TestConstrainedIntEncode(t *testing.T) {
	e := Coder{buf: []byte{0}, isAligned: true}
	ue := Coder{buf: []byte{0}, isAligned: false}
	for i, v := range tbConstrNumA {
		e.appendConstrainedInt64(v.V, v.Min, v.Max)
		ue.appendConstrainedInt64(v.V, v.Min, v.Max)
		if !equal(e.buf, v.Ref) {
			t.Errorf("%d: APER Constrained INTEGER(%d..%d): \nWant %08b \nGot  %08b\n\tLength Encodeod: %d (MUSTBE: %d)\n",
				i, v.Min, v.Max, v.Ref, e.buf, e.BitLen(), v.AL)
		}
		if !equal(ue.buf, v.RefU) || ue.BitLen() != v.UL {
			t.Errorf("%d: UPER Constrained INTEGER(%d..%d): \nWant %08b \nGot  %08b\n\tLength Encodeod: %d(MUSTBE: %d)\n",
				i, v.Min, v.Max, v.RefU, ue.buf, ue.BitLen(), v.UL)
		}
	}
}

var tableUnConstInt = []struct {
	V    []int64
	Ref  []byte
	RefU []byte
}{
	{[]int64{234}, []byte{0x02, 0x00, 0xEA}, []byte{0x02, 0x00, 0xEA}},
	{[]int64{-67}, []byte{0x01, 0xBD}, []byte{0x01, 0xBD}},
	{[]int64{17, 256, 18976}, []byte{0x01, 0x11, 0x02, 0x01, 0x00, 0x02, 0x4A, 0x20}, []byte{0x01, 0x11, 0x02, 0x01, 0x00, 0x02, 0x4A, 0x20}}}

func TestUnConstrinedIntEncode(t *testing.T) {
	for i, v := range tableUnConstInt {
		e := Coder{buf: []byte{0}, isAligned: true}
		ue := Coder{buf: []byte{0}, isAligned: false}
		for _, num := range v.V {
			e.appendUnconstrainedInt64(num)
			ue.appendUnconstrainedInt64(num)
		}
		if !equal(e.buf, v.Ref) {
			t.Errorf("%d: APER Unconstrained Want %08b \nGot  %08b\n",
				i, v.RefU, e.buf)
		}
	}
}

func BenchmarkConstrainedIntEncodeA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := Coder{buf: []byte{0}, isAligned: true}
		for _, v := range tbConstrNumA {
			e.appendConstrainedInt64(v.V, v.Min, v.Max)
		}
	}
}

func BenchmarkConstrainedIntEncodeU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := Coder{buf: []byte{0}, isAligned: false}
		for _, v := range tbConstrNumA {
			e.appendConstrainedInt64(v.V, v.Min, v.Max)
		}
	}
}
