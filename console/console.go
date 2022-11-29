package console

import (
	"Gofreeproxy/fofa"
	"Gofreeproxy/queue"
	"bufio"
	"fmt"
	"github.com/gookit/color"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var liveres []string

func changesocks(ws *net.TCPConn) {
	socksproxy := liveres[rand.Intn(len(liveres))]
	color.RGBStyleFromString("249,134,134").Printf(fmt.Sprintf("\u001B[2K\r[+]当前使用代理%s", socksproxy))
	defer ws.Close()
	//socks, err := net.Dial("tcp", "221.217.53.107:1080")
	socks, err := net.DialTimeout("tcp", socksproxy, 5*time.Second)
	//socks, err := net.Dial("tcp", socksproxy)
	if err != nil {
		log.Println("dial socks error:", err)
		for i := 0; i < len(liveres); i++ {
			if liveres[i] == socksproxy {
				liveres = append(liveres[:i], liveres[i+1:]...)
			}
		}
		changesocks(ws)
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

func Strartsocks() {
	listener, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		log.Fatal(err)
	}
	optSys := runtime.GOOS
	if strings.Contains(optSys, "linux") {
		//执行clear指令清除控制台
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			log.Println("cmd:", err)
		}
	} else {
		//执行clear指令清除控制台
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			log.Println("cmd:", err)
		}
	}
	color.RGBStyleFromString("237,64,35").Printf("[+]一共获取存活代理:%d条\r\n", len(liveres))
	color.RGBStyleFromString("237,64,35").Println("[+]开始监听socks端口: 127.0.0.1:1080")

	for {
		conn, _ := listener.Accept()
		go changesocks(conn.(*net.TCPConn))
	}
}
func IsProxy(proxyIp string, Time int) (isProxy bool) {
	proxyUrl := fmt.Sprintf("socks5://%s", proxyIp)
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return false
	}
	netTransport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	client := &http.Client{
		Timeout:   time.Duration(Time) * time.Second, //设置连接超时时间
		Transport: netTransport,
	}

	res, err := client.Get("http://myip.ipip.net")
	if err != nil {
		return false
	} else {
		defer res.Body.Close()
		if res.StatusCode == 200 {
			body, err := ioutil.ReadAll(res.Body)
			if err == nil && strings.Contains(string(body), "当前 IP") {
				fmt.Printf("\u001B[2K\r[+]%s", string(body))
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
}

func Startgetsocks(Coroutine int, Time int) {
	// protocol="socks5" && "Method:No Authentication(0x00)" && country!="CN" 
	// fofa 排除国内IP
	// keys := "protocol=\"socks5\" && \"Method:No Authentication(0x00)\""
	keys := "protocol=\"socks5\" && \"Method:No Authentication(0x00)\" && country!=\"CN\" "
	GETRES := fofa.Fafaall(keys)
	color.RGBStyleFromString("237,64,35").Printf("[+]从fofa获取代理:%d条", len(GETRES))
	color.RGBStyleFromString("244,211,49").Println("\r\n[+]开始存活性检测")
	pool := queue.New(Coroutine)
	currentdata := 0
	tempsocks := ""
	fmt.Print("\033[s")
	for i := 0; i < len(GETRES); i++ {
		pool.Add(1)
		tempsocks = GETRES[i]
		go func(tempsocks string) {
			flag := IsProxy(tempsocks, Time)
			if flag == true {
				liveres = append(liveres, tempsocks)
			}
			currentdata = currentdata + 1
			pool.Done()
			fmt.Printf("\u001B[2K\r[+]已检测%.2f%%,当前检测IP为:%s", float32(currentdata*100)/float32(len(GETRES)), tempsocks)
		}(tempsocks)

	}

	pool.Wait()
	Writeproxytxt(liveres)
}

func Readfileproxy(Coroutine int, Time int) {
	var fileproxy []string
	fi, err := os.Open("proxy.txt")
	if err != nil {
		log.Println(err)
	}

	// 创建 Reader
	r := bufio.NewReader(fi)
	for {
		line, err := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil && err != io.EOF {
			log.Println(err)
		}
		if err == io.EOF {
			break
		}
		fileproxy = append(fileproxy, line)
	}
	pool := queue.New(Coroutine)
	currentdata := 0
	tempsocks := ""
	fmt.Print("\033[s")
	for i := 0; i < len(fileproxy); i++ {
		pool.Add(1)
		tempsocks = fileproxy[i]
		go func(tempsocks string) {
			//lag := sockslivecheck(tempsocks, client, req)
			flag := IsProxy(tempsocks, Time)
			if flag == true {
				liveres = append(liveres, tempsocks)
			}
			currentdata = currentdata + 1
			pool.Done()
			fmt.Printf("\u001B[2K\r[+]已检测%.2f%%，%s", float32(currentdata*100)/float32(len(fileproxy)), "当前检测IP:"+tempsocks)
		}(tempsocks)

	}

	pool.Wait()
	color.RGBStyleFromString("237,64,35").Printf("[+]一共获取存活代理:%d条", len(liveres))
	fmt.Println(liveres)
}
func Writeproxytxt(livesocks []string) (flag bool) {
	// 指定模式打开文件  追加 文件不存在则创建
	file, err := os.OpenFile("proxy.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// 打开异常检测
	if err != nil {
		fmt.Printf("open file failed, err: %v\n", err)
		return false
	}

	// 延后关闭文件
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Printf("close file failed, err: %v\n", err)
			return
		}
	}(file)
	// 方式一 以二进制形式 写入数据到文件
	for i := 0; i < len(livesocks); i++ {
		_, err = file.Write([]byte(livesocks[i] + "\n"))
		if err != nil {
			return false
		}
	}
	return true
}
