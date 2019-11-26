package asn1

import (
	"testing"
)

var rangeParseT = []struct{
	input	string
	result  []*IntRange
}{
	{ "(4000505005|7888888)", []*IntRange{NewSingleRange(4000505005), NewSingleRange(7888888)}},
	{ "(-17..28|40..50|6)", []*IntRange{NewIntRange(-17,28), NewIntRange(40,50) ,NewSingleRange(6)}},
	{ "(20..35)", []*IntRange{NewIntRange(20,35)}}}


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
