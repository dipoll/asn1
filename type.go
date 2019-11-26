package asn1

type encoder interface {
	Encode() int
}

// Validator validates value against contraints
type Validator interface {
	Validate() error
}
