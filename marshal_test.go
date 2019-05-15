package asn1per

import (
	"fmt"
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

func TestIntegerEncode(t *testing.T) {
	for i, v := range tbBoolNum {
		e := Coder{}
		e.addBool(v.Boolean1)
		e.addBool(v.Boolean1)
		e.appendUint64(7, 6)
		e.appendUint64(3, 3)
		fmt.Printf("Test %d - %X\n", i, e.buf)
		printBytes(e.buf)
		fmt.Printf("Offset: %d\n", e.offset)
	}
	//e.addBool(false)

	fmt.Printf("Should be: %08b %08b\n", 0xC7, 0x60)
	//fmt.Printf("BIN PER: %b\n", boolNumPER)
	//fmt.Printf("BIN UPER: %b\n", boolNumUPER)
}

func TestConstrainedIntUEncode(t *testing.T) {

}

var tbConstrNumA = []struct {
	V   int64
	Min int64
	Max int64
	Ref []byte
}{
	{-3, -5, 0, []byte{0x40}},
	{127, 0, 255, []byte{0x40, 0x7F}},
	{256, 0, 256, []byte{0x40, 0x7F, 0x01, 0x00}}}

func TestConstrainedIntAEncode(t *testing.T) {
	e := Coder{buf: []byte{0}, isAligned: true}
	for i, v := range tbConstrNumA {
		e.appendConstrainedUint64(v.V, v.Min, v.Max)
		if !equal(e.buf, v.Ref) {
			t.Errorf("%d: APER Constrained INTEGER(%d..%d): Want %064b \nGot  %064b\n",
				i, v.Min, v.Max, v.Ref, e.buf)
		}
	}
}
