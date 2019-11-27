package spec

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
