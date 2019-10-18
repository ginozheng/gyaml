package charset

type Codec interface {
	Encode(charSequence []rune) []byte
	Decode(buffer []byte) []rune
}
