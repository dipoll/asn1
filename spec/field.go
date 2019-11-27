package spec

import (
	"encoding/asn1"
	"strconv"
	"strings"
)

type FieldSpec struct {
	Optional     bool   // true iff the field is OPTIONAL
	Explicit     bool   // true iff an EXPLICIT tag is in use.
	Application  bool   // true iff an APPLICATION tag is in use.
	Private      bool   // true iff a PRIVATE tag is in use.
	DefaultValue *int64 // a default value for INTEGER typed fields (maybe nil).
	Tag          *int   // the EXPLICIT or IMPLICIT tag (maybe nil).
	StringType   int    // the string tag to use when marshaling.
	TimeType     int    // the time tag to use when marshaling.
	Set          bool   // true iff this should be encoded as a SET
	OmitEmpty    bool   // true iff this should be omitted if empty when marshaling.
	Range        []IntRange
}

// ParseFieldParameters Given a tag string with the format specified in the package comment,
// parseFieldParameters will parse it into a fieldParameters structure,
// ignoring unknown parts of the string.
func ParseFieldSpec(str string) (ret FieldSpec) {
	for _, part := range strings.Split(str, ",") {
		switch {
		case part == "optional":
			ret.Optional = true
		case part == "explicit":
			ret.Explicit = true
			if ret.Tag == nil {
				ret.Tag = new(int)
			}
		case part == "generalized":
			ret.TimeType = asn1.TagGeneralizedTime
		case part == "utc":
			ret.TimeType = asn1.TagUTCTime
		case part == "ia5":
			ret.StringType = asn1.TagIA5String
		case part == "printable":
			ret.StringType = asn1.TagPrintableString
		case part == "numeric":
			ret.StringType = asn1.TagNumericString
		case part == "utf8":
			ret.StringType = asn1.TagUTF8String
		case strings.HasPrefix(part, "default:"):
			i, err := strconv.ParseInt(part[8:], 10, 64)
			if err == nil {
				ret.DefaultValue = new(int64)
				*ret.DefaultValue = i
			}
		case strings.HasPrefix(part, "tag:"):
			i, err := strconv.Atoi(part[4:])
			if err == nil {
				ret.Tag = new(int)
				*ret.Tag = i
			}
		case part == "set":
			ret.Set = true
		case part == "application":
			ret.Application = true
			if ret.Tag == nil {
				ret.Tag = new(int)
			}
		case part == "private":
			ret.Private = true
			if ret.Tag == nil {
				ret.Tag = new(int)
			}
		case part == "omitempty":
			ret.OmitEmpty = true
		}
	}
	return
}
