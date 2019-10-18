package charset

const bom_UTF8 = "\xef\xbb\xbf"

// UTF-8 字符集
// 可以根据需要自行设置是否过滤BOM头
type UTF8 struct {
	NeedFilterBom bool
	hadFilterBom  bool
}

// 编码
// 将输入的字符序列转换成相应的字节数组
// 可多次调用，不会生成BOM头
func (c *UTF8) Encode(charSequence []rune) []byte {
	return nil
}

// 解码
// 将输入的字节数组转换成相应的字节序列
// 可多次调用，首次调用会根据设置选择进行BOM头的过滤
func (c *UTF8) Decode(buffer []byte) []rune {
	pos, length := 0, len(buffer)
	charSequence := make([]rune, 0, length/2)
	if c.NeedFilterBom && !c.hadFilterBom {
		if length < 3 {
			return charSequence
		} else {
			if buffer[0] == bom_UTF8[0] && buffer[1] == bom_UTF8[1] && buffer[2] == bom_UTF8[2] {
				pos += 3
			}
			c.hadFilterBom = true
		}
	}

	for pos != length {

	}
	return nil
}
