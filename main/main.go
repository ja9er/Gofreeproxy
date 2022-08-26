package main

import (
	"Gofreeproxy/fofa"
	"Gofreeproxy/queue"
	"crypto/tls"
	"fmt"
	"github.com/gookit/color"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func changesocks(ws *net.TCPConn) {
	defer ws.Close()
	//socks, err := net.Dial("tcp", "221.217.53.107:1080")
	socks, err := net.Dial("tcp", "47.74.70.193:45554")
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
	url_i := url.URL{}
	url_proxy, _ := url_i.Parse("socks://" + SocksProxy)
	tr := &http.Transport{
		//关闭证书验证
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//设置超时
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		//设置代理
		Proxy: http.ProxyURL(url_proxy),
	}
	client.Transport = tr
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	httpbody := string(result)
	if len(httpbody) > 10 && resp.StatusCode == 200 {
		return true
	}
	return false

}

func startgetsocks() {
	keys := "protocol=\"socks5\" && \"Method:No Authentication(0x00)\" && port=\"1080\""
	GETRES := fofa.Fafaall(keys)
	color.RGBStyleFromString("237,64,35").Printf("[+]从fofa获取代理:%d条", len(GETRES))
	color.RGBStyleFromString("244,211,49").Println("\r\n[+]开始存活性检测")
	pool := queue.New(100)
	var liveres []string
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
	for i := 0; i < len(GETRES); i++ {
		pool.Add(1)
		tempsocks := GETRES[i]
		go func(tempsocks string) {
			flag := sockslivecheck(tempsocks, client, req)
			if flag == true {
				fmt.Printf("[+]已检测%.2f%%,存活检测成功:%s\r", float32(i)/float32(len(GETRES)), GETRES[i-1])
				liveres = append(liveres, GETRES[i-1])
			} else {
				fmt.Printf("[-]已检测%.2f%%,存活检测失败:%s\r", float32(i)/float32(len(GETRES)), GETRES[i-1])
			}
			pool.Done()
		}(tempsocks)
	}
	pool.Wait()
	color.RGBStyleFromString("237,64,35").Printf("[+]一共获取存活代理:%d条", len(liveres))
	fmt.Println(liveres)
}
func main() {
	startgetsocks()
	//strartsocks()
	//falg:=sockslivecheck("171.214.198.59:1080")
	//fmt.Println(falg)
}
