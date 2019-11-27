package asn1

import (
	"encoding/asn1"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ParseFieldParameters Given a tag string with the format specified in the package comment,
// parseFieldParameters will parse it into a fieldParameters structure,
// ignoring unknown parts of the string.
func ParseFieldParameters(str string) (ret FieldParameters) {
	for _, part := range strings.Split(str, ",") {
		switch {
		case part == "optional":
			ret.optional = true
		case part == "explicit":
			ret.explicit = true
			if ret.tag == nil {
				ret.tag = new(int)
			}
		case part == "generalized":
			ret.timeType = asn1.TagGeneralizedTime
		case part == "utc":
			ret.timeType = asn1.TagUTCTime
		case part == "ia5":
			ret.stringType = asn1.TagIA5String
		case part == "printable":
			ret.stringType = asn1.TagPrintableString
		case part == "numeric":
			ret.stringType = asn1.TagNumericString
		case part == "utf8":
			ret.stringType = asn1.TagUTF8String
		case strings.HasPrefix(part, "default:"):
			i, err := strconv.ParseInt(part[8:], 10, 64)
			if err == nil {
				ret.defaultValue = new(int64)
				*ret.defaultValue = i
			}
		case strings.HasPrefix(part, "tag:"):
			i, err := strconv.Atoi(part[4:])
			if err == nil {
				ret.tag = new(int)
				*ret.tag = i
			}
		case part == "set":
			ret.set = true
		case part == "application":
			ret.application = true
			if ret.tag == nil {
				ret.tag = new(int)
			}
		case part == "private":
			ret.private = true
			if ret.tag == nil {
				ret.tag = new(int)
			}
		case part == "omitempty":
			ret.omitEmpty = true
		}
	}
	return
}

// GetTagParams returns ASN1 tag and parsed params
// this is
func GetTagParams(value interface{}) (int, *FieldParameters) {
	v := reflect.TypeOf(value)

	switch v.Kind() {
	case reflect.Struct:
		fmt.Println("This value is a struct: name: ", v.String())
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			fmt.Println("Name: ", f.Name)
			fmt.Println("Tag is:", f.Tag.Get("asn1"))
		}
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return asn1.TagInteger, &FieldParameters{}
	}
	return 0, nil
}
