package asn1per

type encoder interface {
	Encode() int
}

// Validator validates value against contraints
type Validator interface {
	Validate() error
}
