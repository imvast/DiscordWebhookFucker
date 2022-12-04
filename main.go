package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/admin100/util/console"
	"github.com/corpix/uarand"
	"github.com/gosuri/uilive"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"math/rand"
	"os"
	"sync"
	"time"
	"main/utils"
)


var (
	Config = utils.LoadConfig()

	sentamt  int
	workdamt int

	start int
	end   int

	proxies []string

	red    = "\x1b[38;5;167m"
	purple = "\x1b[38;5;57m"
	pink   = "\x1b[38;5;128m"
	reset  = "\x1b[0m"
)


func main() {
	prox, err := os.Open("./proxies.txt")
	if err != nil {
		fmt.Printf("[%s%s%s] Failed to read ./proxies.txt\n", red, time.Now().Format("15:04:05"), reset)
		return
	}
	s := bufio.NewScanner(prox)
	for s.Scan() {
		proxies = append(proxies, s.Text())
	}

	go TitleThread()
	go ConsoleThread()

	start = int(time.Now().Unix())
	wg := sync.WaitGroup{}
	goroutines := make(chan struct{}, Config.Threads)
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			goroutines <- struct{}{}
			SendReq()
			<-goroutines
		}()
	}
	wg.Wait()
}

func ConsoleThread() {
	fmt.Printf(purple+`
			▓█████▄  █     █░  █████▒
			▒██▀ ██▌▓█░ █ ░█░▓██   ▒ 
			░██   █▌▒█░ █ ░█ ▒████ ░ 
			░▓█▄   ▌░█░ █ ░█ ░▓█▒  ░ 
			░▒████▓ ░░██▒██▓ ░▒█░    
			▒▒▓  ▒ ░ ▓░▒ ▒   ▒ ░    
			░ ▒  ▒   ▒ ░ ░   ░      
			░ ░  ░   ░   ░   ░ ░    
			░        ░            
			░                       
    %v[%vDiscord Webhook Fucker%v]%v               %v$%vimvast%v

		           %vProxies Loaded:%v %v
	%v`, pink, purple, pink, reset, pink, purple, reset, purple, pink, len(proxies), reset)
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()
	for {
		fmt.Fprintf(writer, "%vWebhooks  Sent:%v %v\n", purple, pink, workdamt)
		fmt.Fprintf(writer.Newline(), "\n%v[%v] %v[Running Fucker]%v", pink, time.Now().Format("15:04:05"), purple, reset)
		time.Sleep(time.Millisecond * 5)
		writer.Flush()
	}
}

func TitleThread() {
	for {
		console.SetConsoleTitle(fmt.Sprintf("imvast~DWF | Sent: %d/%d - Elapsed: %ds", workdamt, sentamt, (int(time.Now().Unix()) - start)))
		time.Sleep(500)
	}
}

func SendReq() {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	values := map[string]string{
		"content":    Config.Content,
		"username":   "github.com/imvast",
		"avatar_url": Config.AvatarUrl}
	jsonValue, _ := json.Marshal(values)

	req.Header.SetMethod("POST")
	req.SetRequestURI(Config.Webhook)
	req.SetBody(jsonValue)
	req.Header.Set("user-agent", uarand.GetRandom())
	req.Header.SetContentType("application/json")

	proxy := proxies[rand.Intn(len(proxies))]
	client := &fasthttp.Client{
		Dial:           fasthttpproxy.FasthttpHTTPDialer(proxy),
		ReadBufferSize: 50_000,
	}

	err := client.Do(req, res)
	if err != nil {
		return //fmt.Printf(err.Error())
	}
	sentamt++

	if res.StatusCode() == 204 {
		workdamt++
	} else if res.StatusCode() == 429 && Config.Debug == true{
		body := string(res.Body())
		fmt.Printf("[!] %v | %v", res.StatusCode(), body)
	} else {
		return //fmt.Println("check for error: res.StatusCode")
	}
}
