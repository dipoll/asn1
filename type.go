package asn1per


type encoder interface {
	Encode() int
}


type Validator interface {
	Validate() error
} 

type IntRange struct {
	Min		int
	Max		int
}
