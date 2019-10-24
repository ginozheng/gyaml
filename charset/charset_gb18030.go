package charset

type GB18030 struct{}

func (c *GB18030) Encode(runes []rune) ([]byte, error) {
	if len(runes) == 0 {
		return nil, &EncodingError{msg: "Zero rune set length."}
	}

	bytes := make([]byte, 0)
	for _, v := range runes {
		value := uint32(v)
		if value >= 0xd800 && value <= 0xdfff {
			return nil, &EncodingError{msg: "GB18030 doesn't support number between 0xd800 and 0xdfff."}
		}

		if value <= 0x80 {
			bytes = append(bytes, byte(value))
		} else if (value >= 0x8140 && value <= 0xfe7e) || (value >= 0x8180 && value <= 0xfefe) {
			firstOctet := byte(value >> 8)
			bytes = append(bytes, firstOctet)

			secondOctet := byte(value & 0xff)
			bytes = append(bytes, secondOctet)
		} else if value >= 0x81308130 && value <= 0xFE39FE39 {
			firstOctet := byte(value >> 24)
			bytes = append(bytes, firstOctet)

			secondOctet := byte(value >> 16 & 0xff)
			bytes = append(bytes, secondOctet)

			thirdOctet := byte(value >> 8 & 0xff)
			bytes = append(bytes, thirdOctet)

			fourthOctet := byte(value & 0xff)
			bytes = append(bytes, fourthOctet)
		} else {
			return nil, &EncodingError{msg: "Invalid number."}
		}
	}

	return bytes, nil
}

func (c *GB18030) Decode(bytes []byte) ([]rune, error) {
	pos, length := 0, len(bytes)

	if length == 0 {
		return nil, &DecodingError{msg: "Zero byte set length."}
	}

	runes := make([]rune, 0)
	for pos != length {
		var value rune
		firstOctet := bytes[pos]

		if firstOctet == 0xff {
			return nil, &DecodingError{msg: "GB18030 doesn't support first byte 0xff."}
		}

		if firstOctet <= 0x80 {
			// first byte 0x00~0x80
			value = rune(firstOctet)
			runes = append(runes, value)
			pos += 1
		} else {
			// first byte 0x81~0xfe
			if pos+1 >= length {
				return nil, &DecodingError{msg: "No enough bytes for decoding."}
			}

			secondOctet := bytes[pos+1]
			if secondOctet >= 0x40 && secondOctet <= 0xfe {
				if secondOctet == 0x7f {
					return nil, &DecodingError{msg: "Second octet doesn't support 0x7f which width is 2."}
				}

				value = rune(firstOctet<<8 | secondOctet)
				runes = append(runes, value)
				pos += 2
			} else if secondOctet >= 0x30 && secondOctet <= 0x39 {
				if pos+3 >= length {
					return nil, &DecodingError{msg: "No enough bytes for decoding which width is 4."}
				}

				thirdOctet := bytes[pos+2]
				if !(thirdOctet >= 0x81 && thirdOctet <= 0xfe) {
					return nil, &DecodingError{msg: "Third octet should between 0x81 and 0xfe which width is 4."}
				}

				fourthOctet := bytes[pos+3]
				if !(fourthOctet >= 0x30 && fourthOctet <= 0x39) {
					return nil, &DecodingError{msg: "Fourth octet should between 0x30 and 0x39 which width is 4."}
				}

				value = rune(firstOctet<<24 | secondOctet<<16 | thirdOctet<<8 | fourthOctet)
				if value >= 0xd800 && value <= 0xdfff {
					return nil, &DecodingError{msg: "GB18030 doesn't support number between 0xd800 and 0xdfff"}
				}

				runes = append(runes, value)
				pos += 4
			} else {
				return nil, &DecodingError{msg: "Invalid byte set."}
			}
		}
	}
}
