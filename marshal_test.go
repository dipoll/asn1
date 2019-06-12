package asn1per

import (
	"testing"
)


func TestParseRange(t *testing.T){
	rng, err := ParseRange("(4000505005|7888888)")
	if len(rng) != 2 || err != nil {
		t.Errorf("bad parsing of range|size numbers")
	}
	if *rng[0].Max != 4000505005 {
		t.Errorf("bad parsing of number 4000505005")
	}
	if *rng[1].Max != 7888888 {
			t.Errorf("bad parsing of number 7888888")
	}
}