package main

import (
	"fmt"
	"github.com/ginozheng/gyaml/charset"
	"io/ioutil"
)

func main() {
	cc := charset.UTF8{}
	a := "我是中国人"
	runes, er := cc.Decode([]byte(a))
	if er != nil {
		fmt.Println(er.Error())
	}

	if runes[0] == '我' {
		fmt.Println("yes")
	}

	for _, v := range runes {
		fmt.Printf("%#U\n", v)
	}

	fmt.Println("====================================")
	bytes, er1 := cc.Encode([]rune(a))
	if er1 != nil {
		fmt.Println(er1.Error())
	}
	for _, v := range bytes {
		fmt.Println(v)
	}

	fmt.Println("====================================")
	bs := []byte(a)
	for _, v := range bs {
		fmt.Println(v)
	}

	bss, err := ioutil.ReadFile("test.txt")
	if err != nil {
		fmt.Println(err.Error())
	}
	gb := charset.GB18030{}
	runes, er = gb.Decode(bss)
	if er != nil {
		fmt.Println(er.Error())
	}
	for _, v := range runes {
		fmt.Printf("%#U\n", v)
	}
}
