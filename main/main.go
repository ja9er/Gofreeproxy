package main

import (
	"Gofreeproxy/fofa"
	"Gofreeproxy/queue"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"github.com/gookit/color"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

func changesocks(ws *net.TCPConn) {
	defer ws.Close()
	//socks, err := net.Dial("tcp", "221.217.53.107:1080")
	socks, err := net.Dial("tcp", "60.190.195.146:1080")
	if err != nil {
		log.Println("dial socks error:", err)
		return
	}
	defer socks.Close()
	var wg sync.WaitGroup
	ioCopy := func(dst io.Writer, src io.Reader) {
		defer wg.Done()
		io.Copy(dst, src)
	}
	wg.Add(2)
	go ioCopy(socks, ws)
	go ioCopy(ws, socks)
	wg.Wait()
}

func strartsocks() {
	listener, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listen tcp at:", "127.0.0.1:1080")
	for {
		conn, _ := listener.Accept()
		go changesocks(conn.(*net.TCPConn))
	}
}
func sockslivecheck(SocksProxy string, client *http.Client, req *http.Request) bool {
	dialer, err := proxy.SOCKS5("tcp", SocksProxy, nil, proxy.Direct)
	if err != nil {
		log.Println("can't connect to the proxy:", err)
		_, err2 := proxy.SOCKS5("udp", SocksProxy, nil, proxy.Direct)
		if err2 == nil {
			log.Println("[+]udp proxy:", err)
			return false
		}
		return false
	}
	tr := &http.Transport{
		//关闭证书验证
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//设置超时
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		//设置代理

	}
	tr.Dial = dialer.Dial
	client.Transport = tr
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return false
		}
	}
	result, err := ioutil.ReadAll(reader)
	httpbody := string(result)
	if strings.Contains(httpbody, "当前 IP") {
		fmt.Printf("[+]%s,使用代理:%s\r", httpbody, SocksProxy)
		return true
	}
	return false

}
func startgetsocks() {
	keys := "protocol=\"socks5\" && \"Method:No Authentication(0x00)\" && port=\"1080\""
	GETRES := fofa.Fafaall(keys)
	color.RGBStyleFromString("237,64,35").Printf("[+]从fofa获取代理:%d条", len(GETRES))
	color.RGBStyleFromString("244,211,49").Println("\r\n[+]开始存活性检测")
	pool := queue.New(200)
	var liveres []string
	client := &http.Client{
		//禁止重定向
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, _ := http.NewRequest("GET", "http://myip.ipip.net", nil)
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.3; rv:36.0) Gecko/20100101 Firefox/36.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	currentdata := 0
	tempsocks := ""
	for i := 0; i < len(GETRES); i++ {
		pool.Add(1)
		tempsocks = GETRES[i]
		go func(tempsocks string) {
			flag := sockslivecheck(tempsocks, client, req)
			if flag == true {
				liveres = append(liveres, tempsocks)
			}
			currentdata = currentdata + 1
			pool.Done()
			fmt.Printf("[+]已检测%.2f%%,当前检测IP:%s,当前检测完成总数：%d\r", float32(currentdata*100)/float32(len(GETRES)), tempsocks, currentdata)
		}(tempsocks)
	}

	pool.Wait()
	color.RGBStyleFromString("237,64,35").Printf("[+]一共获取存活代理:%d条", len(liveres))
	fmt.Println(liveres)
}

func test(socksproxy string) {
	client := &http.Client{
		//禁止重定向
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, _ := http.NewRequest("GET", "https://api.ip.sb/ip", nil)
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.3; rv:36.0) Gecko/20100101 Firefox/36.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	flag := sockslivecheck(socksproxy, client, req)
	fmt.Println(flag)
}

func main() {
	startgetsocks()
	//strartsocks()
	//test("127.0.0.1:7890")

}
