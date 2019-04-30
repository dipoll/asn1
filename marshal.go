package asn1per

import (
	"encoding/asn1"
	"fmt"
)

func parseBool(bytes []byte) (bool, error) {
	fmt.Println(bytes[0], bytes[0]>>7)
	if len(bytes) != 1 {
		return false, asn1.StructuralError{"invalid boolean"}
	}
	switch bytes[0] >> 7 {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, asn1.SyntaxError{"invalid boolean value"}
	}
	return false, nil
}
