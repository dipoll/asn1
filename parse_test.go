package asn1

import (
	"fmt"
	"math/big"
	"testing"
)

type myTestData struct {
	OrdinaryNumber int `asn1:"range(10..5),integer,something_else"`
	Fixed64Number  int64
	BigNumber      *big.Int
}

func TestStructParse(t *testing.T) {
	var tS myTestData

	tag, fParams := GetTagParams(tS)
	fmt.Println(tag, fParams)
}
