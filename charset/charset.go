package charset

type Codec interface {
	Encode(runes []rune) ([]byte, error)
	Decode(bytes []byte) ([]rune, error)
}

type EncodingError struct {
	msg string
}

func (ee *EncodingError) Error() string {
	return ee.msg
}

type DecodingError struct {
	msg string
}

func (de *DecodingError) Error() string {
	return de.msg
}
