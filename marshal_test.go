package asn1per

import (
	"fmt"
	"testing"
)

func TestBooleanParsing(t *testing.T) {
	fmt.Println(appendBit(byte(0), 0))
	b, err := parseBool([]byte{0x80})
	if err == nil {
		t.Errorf("Failed Parse Boolean: %v %v\n", err, 0x80)
	}
	if b != true {
		t.Errorf("Failed Parse True Boolean Value\n")
	}
	fmt.Printf("%X\n", appendBit(byte(0), 1))
}
