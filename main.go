package main

import (
	"fmt"
	"github.com/next-framework/nextyaml/charset"
)

func main() {
	cc := charset.UTF8{}
	a := "我是中国人"
	er, runes := cc.Decode([]byte(a))
	if er != nil {
		fmt.Println(er.Error())
	}
	for _, v := range runes {
		fmt.Printf("%#U\n", v)
	}

	fmt.Println("====================================")
	er1, bytes := cc.Encode([]rune(a))
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
}
