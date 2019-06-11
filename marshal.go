package asn1per

import (
	"errors"
	"strconv"
	"strings"
)

type intRange struct {
	Min *int
	Max *int
}

// ParseRange parsers string into
// set of ranges:
//	(1..5) - simple single range
//  (1..5|8..10) - multiple ranges
//  (1|5|25) - choises only with
func ParseRange(s string) ([]intRange, error) {
	s = strings.TrimLeft(s, "(")
	s = strings.TrimRight(s, ")")
	var out []intRange
	values := strings.Split(s, "|")

	for _, v := range values {
		vl := strings.Split(v, "..")
		switch len(vl) {
		case 2:
			v1, _ := strconv.Atoi(vl[0])
			v2, _ := strconv.Atoi(vl[1])
			out = append(out, intRange{&v1, &v2})
		case 1:
			v2, _ := strconv.Atoi(vl[0])
			out = append(out, intRange{Max: &v2})
		default:
			return out, errors.New("asn1: tag: can not parse range or size")
		}
	}

	return out, nil

}
