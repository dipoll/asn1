package asn1per

import (
	"testing"
)

var rangeParseT = []struct{
	input	string
	result  []*IntRange
}{
	{ "(4000505005|7888888)", []*IntRange{NewSingleRange(4000505005), NewSingleRange(7888888)}}}


func TestParseRange(t *testing.T){

	for i, v := range rangeParseT {
		rng, err := ParseRange(v.input)
		if len(rng) != len(v.result) || err != nil {
			t.Errorf("asn1: [%d] bad parsing of range|size numbers", i )
		}
		for num, rValue := range rng {
			if !rValue.Equals(v.result[num]) {
				t.Errorf("asn1: [%d] bad parsing of range|size numbers: expect: %v, got: %v", i,v.result[num] ,rValue )
			}
			
		}
	}
}
