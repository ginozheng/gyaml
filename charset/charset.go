package charset

type Codec interface {
	Encode(buffer []byte) []rune
	Decode(charSequence []rune) []byte
}
