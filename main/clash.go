package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// 代理节点
type nodes struct {
	Nodename   string `yaml:"name"`
	Nodeport   int    `yaml:"port"`
	NodeServer string `yaml:"server"`
	Nodetype   string `yaml:"type"`
}

// 代理组
type proxy_groups struct {
	Name     string   `yaml:"name"`
	Interval int      `yaml:"interval"`
	Node     []string `yaml:"proxies"`
	Strategy string   `yaml:"strategy"`
	Type     string   `yaml:"type"`
	Url      string   `yaml:"url"`
}

// 基础配置
type setconf struct {
	Allowlan           bool           `yaml:"allow-lan"`
	Port               int            `yaml:"port"`
	Socksport          int            `yaml:"socks-port"`
	Redirport          int            `yaml:"redir-port"`
	Mode               string         `yaml:"mode"`
	Loglevel           string         `yaml:"log-level"`
	Externalcontroller string         `yaml:"external-controller"`
	Proxies            []nodes        `yaml:"proxies"`
	Proxygroup         []proxy_groups `yaml:"proxy-groups"`
	Rules              []string       `yaml:"rules"`
}

// 写入yaml
func writeToXml(src string, Nodeall []nodes) {
	fmt.Printf("v: %v\n", Nodeall)
	stu := &setconf{
		Allowlan:           true,
		Port:               7890,
		Socksport:          7891,
		Redirport:          7892,
		Mode:               "global",
		Loglevel:           "info",
		Externalcontroller: "127.0.0.1:9090",
		Proxies:            Nodeall,
		Proxygroup: []proxy_groups{{
			"负载模式", 100, []string{"测试节点1", "测试节点2"},
			"round-robin", "load-balance", "https://www.baidu.com/"}},
		Rules: []string{
			"DOMAIN-SUFFIX,google.com,DIRECT",
			"DOMAIN-KEYWORD,google,DIRECT",
			"DOMAIN,google.com,DIRECT",
			"DOMAIN-SUFFIX,ad.com,REJECT",
			"GEOIP,CN,DIRECT", "MATCH,DIRECT"},
	}
	data, err := yaml.Marshal(stu)
	checkError(err)
	err = ioutil.WriteFile(src, data, 0777)
	checkError(err)
}

func main() {
	var node = []nodes{
		{"测试节点1", 7845, "8.8.8.8", "socks5"},
		{"测试节点2", 7845, "8.8.8.8", "socks5"},
	}
	//node := nodes{"2133",7845,"8.8.8.8","socks5"}
	src := "test1.yaml"
	writeToXml(src, node)
}
