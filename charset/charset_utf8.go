package charset

// UTF-8
// See more detail from https://tools.ietf.org/html/rfc3629
// For simply:
//   The table below summarizes the format of these different octet types.
//   The letter x indicates bits available for encoding bits of the
//   character number.
//
//   Char. number range  |        UTF-8 octet sequence
//      (hexadecimal)    |              (binary)
//   --------------------+---------------------------------------------
//   0000 0000-0000 007F | 0xxxxxxx
//   0000 0080-0000 07FF | 110xxxxx 10xxxxxx
//   0000 0800-0000 FFFF | 1110xxxx 10xxxxxx 10xxxxxx
//   0001 0000-0010 FFFF | 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
//
//   Encoding a character to UTF-8 proceeds as follows:
//
//   1.  Determine the number of octets required from the character number
//       and the first column of the table above.  It is important to note
//       that the rows of the table are mutually exclusive, i.e., there is
//       only one valid way to encode a given character.
//
//   2.  Prepare the high-order bits of the octets as per the second
//       column of the table.
//
//   3.  Fill in the bits marked x from the bits of the character number,
//       expressed in binary.  Start by putting the lowest-order bit of
//       the character number in the lowest-order position of the last
//       octet of the sequence, then put the next higher-order bit of the
//       character number in the next higher-order position of that octet,
//       etc.  When the x bits of the last octet are filled in, move on to
//       the next to last octet, then to the preceding one, etc. until all
//       x bits are filled in.
type UTF8 struct{}

func (c *UTF8) Encode(runes []rune) ([]byte, error) {
	if len(runes) == 0 {
		return nil, &EncodingError{msg: "Zero rune set length."}
	}

	bytes := make([]byte, 0)
	for _, v := range runes {
		if v >= 0xd800 && v <= 0xdffff {
			return nil, &EncodingError{msg: "UTF-8 doesn't support numbers between 0xd800 and 0xdfff."}
		}

		if v > 0x10ffff {
			return nil, &EncodingError{msg: "UTF-8 doesn't support numbers greater than 0x10ffff."}
		}

		var width int
		var octet byte
		switch {
		case v <= 0x7f:
			octet = byte(v & 0x7f)
			width = 1
		case v >= 0x80 && v <= 0x07ff:
			octet = byte(((v >> 6) | 0xc0) | ((v >> 6) & 0x1f))
			width = 2
		case v >= 0x800 && v <= 0xffff:
			octet = byte(((v >> 12) | 0xe0) | ((v >> 12) & 0x0f))
			width = 3
		case v >= 0x10000 && v <= 0x10ffff:
			octet = byte(((v >> 18) | 0xf0) | ((v >> 18) & 0x07))
			width = 7
		default:
			return nil, &EncodingError{msg: "Invalid rune set."}
		}

		bytes = append(bytes, octet)

		for i := width - 2; i >= 0; i-- {
			octetFollowing := byte((v >> (i * 6)) & 0x3f)
			octetFollowing = (octetFollowing | 0x80) | (octetFollowing & 0x3f)
			bytes = append(bytes, octetFollowing)
		}
	}
	return bytes, nil
}

func (c *UTF8) Decode(bytes []byte) ([]rune, error) {
	pos, length := 0, len(bytes)

	if length == 0 {
		return nil, &DecodingError{msg: "Zero byte set length."}
	}

	runes := make([]rune, 0)
	for pos != length {
		var width int
		var value rune
		octet := bytes[pos]
		switch {
		case octet&0x80 == 0x00:
			width = 1
			value = rune(octet & 0x7f)
		case octet&0xe0 == 0xc0:
			width = 2
			value = rune(octet & 0x1f)
		case octet&0xf0 == 0xe0:
			width = 3
			value = rune(octet & 0x0f)
		case octet&0xf8 == 0xf0:
			width = 4
			value = rune(octet & 0x07)
		default:
			return nil, &DecodingError{msg: "Invalid byte set."}
		}

		if pos+width > length {
			return nil, &DecodingError{msg: "Width error."}
		}

		for i := 1; i < width; i++ {
			octetFollowing := bytes[pos+i]
			if octetFollowing&0xc0 != 0x80 {
			}
			value = value<<6 | rune(octetFollowing&0x3f)
		}

		if value >= 0xd800 && value <= 0xdfff {
			return nil, &DecodingError{msg: "UTF-8 doesn't support numbers between 0xd800 and 0xdfff."}
		}

		if value > 0x10ffff {
			return nil, &DecodingError{msg: "UTF-8 doesn't support numbers greater than 0x10ffff."}
		}

		switch {
		case width == 1:
		case width == 2:
			if value < 0x80 {
				return nil, &DecodingError{msg: "Numbers must between 0x80 and 0x7ff while width is 2."}
			}
		case width == 3:
			if value < 0x800 {
				return nil, &DecodingError{msg: "Numbers must between 0x800 and 0xffff while width is 3."}
			}
		case width == 4:
			if value < 0x10000 {
				return nil, &DecodingError{msg: "Number must between 0x10000 and 0x10ffff while width is 4."}
			}
		}

		runes = append(runes, value)
		pos += width
	}
	return runes, nil
}
