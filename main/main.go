package main

import (
	"Gofreeproxy/console"
	"flag"
	"fmt"
)

var (
	Fofa      bool
	Quake     bool
	Hunter    bool
	Renew     bool
	Port      string
	File      bool
	Coroutine int
	Time      int
	logo      = " _________  ____  ____  _______  ____  _\n/  __/  _ \\/  __\\/  __\\/  _ \\  \\//\\  \\//\n| |  | / \\||  \\/||  \\/|| / \\|\\  /  \\  / \n| |_// \\_/||  __/|    /| \\_/|/  \\  / /  \n\\____\\____/\\_/   \\_/\\_\\\\____/__/\\\\/_/   \n                                        "
)

func main() {
	fmt.Println(logo)
	flag.BoolVar(&Fofa, "fofa", false, "\n使用-fofa参数可从fofa收集资产获取公开代理使用")
	flag.BoolVar(&Hunter, "hunter", false, "\n使用-hunter参数可从fofa收集资产获取公开代理使用")
	flag.BoolVar(&Quake, "quake", false, "\n使用-quake参数可从hunter收集资产获取公开代理使用")
	flag.BoolVar(&Renew, "renew", false, "\n当启用-fofa或-quake或-hunter参数时是否对现有proxy.txt进行重写，默认为否")
	flag.BoolVar(&File, "f", false, "\n使用-f参数可读取当前目录下的proxy.txt，获取其中的代理使用")
	flag.StringVar(&Port, "p", "1080", "\n使用-p参数自定义本地服务端口，默认为1080")
	flag.IntVar(&Coroutine, "c", 200, "\n使用-c参数可设置验证代理的协程数量，默认为200")
	flag.IntVar(&Time, "t", 10, "\n使用-t参数可设置验证代理的超时时间，默认为10秒")
	flag.Parse()

	if Fofa || Quake || Hunter {
		console.Startgetsocks(Coroutine, Time, Fofa, Quake, Hunter, Renew)
		console.Strartsocks(Port)
	} else if File == true {
		console.Readfileproxy(Coroutine, Time)
		console.Strartsocks(Port)
	} else {
		flag.Usage()
		fmt.Println("请输入参数")
	}
}
