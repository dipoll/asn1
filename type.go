package asn1per

type PackType uint8


const (
	ALIGNED		PackType = iota
	UNALIGNED
)


type Codec struct {
	pack	PackType
	bits	uint64	
}


func NewEncoder() *Codec{
	return &Codec{}
}

func NewDecoder() *Codec {
	return &Codec{}
}