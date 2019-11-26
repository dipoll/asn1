package asn1

import (
	"encoding/asn1"
	"errors"
	"strconv"
	"strings"
)

// Extension of the ASN1 tags provided in encoding/asn1 module
//

var o asn1.TagInteger

// IntRange represents integer range or size
type IntRange struct {
	Min *int
	Max *int
}

func (r IntRange) Equals(other *IntRange) bool {
	if other.Min != nil && r.Min != nil {
		if *other.Min != *r.Min {
			return false
		}
	} else if other.Min == nil && r.Min == nil {

	} else {
		return false
	}
	if other.Max != nil && r.Max != nil {
		if *other.Max != *r.Max {
			return false
		}
	} else if other.Max == nil && r.Max == nil {
		return true
	} else {
		return false
	}

	return true

}

// NewIntRange
func NewIntRange(min, max int) *IntRange {
	return &IntRange{&min, &max}
}

// NewIntRange
func NewSingleRange(value int) *IntRange {
	return &IntRange{nil, &value}
}

// ParseRange parsers string into
// set of ranges:
//	(1..5) - simple single range
//  (1..5|8..10) - multiple ranges
//  (1|5|25) - choises only with
func ParseRange(s string) ([]IntRange, error) {
	s = strings.TrimLeft(s, "(")
	s = strings.TrimRight(s, ")")
	var out []IntRange
	values := strings.Split(s, "|")

	for _, v := range values {
		vl := strings.Split(v, "..")
		switch len(vl) {
		case 2:
			v1, _ := strconv.Atoi(vl[0])
			v2, _ := strconv.Atoi(vl[1])
			out = append(out, IntRange{&v1, &v2})
		case 1:
			v2, _ := strconv.Atoi(vl[0])
			out = append(out, IntRange{Max: &v2})
		default:
			return out, errors.New("asn1: tag: can not parse range or size")
		}
	}

	return out, nil

}
