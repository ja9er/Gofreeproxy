package main

import (
	"Gofreeproxy/console"
	"flag"
	"fmt"
)

var (
	Fofa bool
	File bool
	logo = " _________  ____  ____  _______  ____  _\n/  __/  _ \\/  __\\/  __\\/  _ \\  \\//\\  \\//\n| |  | / \\||  \\/||  \\/|| / \\|\\  /  \\  / \n| |_// \\_/||  __/|    /| \\_/|/  \\  / /  \n\\____\\____/\\_/   \\_/\\_\\\\____/__/\\\\/_/   \n                                        "
)

func main() {
	fmt.Println(logo)
	flag.BoolVar(&Fofa, "fofa", false, "使用-fofa参数可从fofa收集资产获取公开代理使用")
	flag.BoolVar(&File, "f", false, "使用-f参数可读取当前目录下的proxy.txt，获取其中的代理使用")
	flag.Parse()

	if Fofa == true {
		console.Startgetsocks()
		console.Strartsocks()
	} else if File == true {
		console.Readfileproxy()
		console.Strartsocks()
	} else {
		flag.Usage()
		//console.Test("39.105.138.223:1080")
		fmt.Println("请输入参数")
	}
	//startgetsocks()
	//strartsocks()
	//test("127.0.0.1:7890")

}
