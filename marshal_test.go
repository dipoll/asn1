package asn1per

import (
	"fmt"
	"testing"
)

var (
	testBool byte = 0xBF // PER and UPER in one byte {1 0 1 1 1 1 1 1}
)

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
	e := Encoder{}
	e.addUnsignedNumber(2)
	fmt.Println("Finished")
}
