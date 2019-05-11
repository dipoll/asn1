package asn1per

import (
	"fmt"
	"testing"
)

var (
	testBool byte = 0xBF // PER and UPER in one byte {1 0 1 1 1 1 1 1}
	boolNumPER []byte = []byte{0x80, 0x01, 0x02} // True and Unconstrained Number
	boolNumUPER []byte = []byte{0x80, 0x81, 0x00}
	boolNumConstPER []byte = []byte{}
)

func equal(a, b []int) bool {
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
		fmt.Printf("POS: %d, Out Value: %X\n", npos, out)
		npos = appendBool(&out, npos, v)
	}
	if out != testBool {
		t.Errorf("Bool Encode Failed: Expect: %X, Got: %X", testBool, out)
	}
}

func TestIntegerEncode(t *testing.T) {
	e := Coder{}
	e.addBool(true)
	e.addUint64(2, 7)
	fmt.Printf("%X\n",e.buf)
}
