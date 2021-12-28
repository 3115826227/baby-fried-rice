package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	addr       = "127.0.0.1:18070"
	userNumber = 50000
	tokens     []string
)

func Post(url string, data string) (body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data))

	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func Get(url string, token string, params map[string]interface{}) (body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("token", token)
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func Websocket(url string, token string) error {
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 60 * time.Second,
	}
	conn, _, err := dialer.Dial(url+"?token="+token, nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			var mp = map[string]interface{}{
				"ws_message_notify_type": 0,
			}
			if err = conn.WriteJSON(mp); err != nil {
				return err
			}
		}
	}
}

type ReqRegister struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
	Username  string `json:"username"`
}

func register(number int) {
	var req = ReqRegister{
		LoginName: fmt.Sprintf("baby_test_%v", number),
		Password:  "123456",
		Username:  fmt.Sprintf("测试账号%v", number),
	}
	data, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = Post(fmt.Sprintf("http://%v/api/user/register", addr), string(data))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

func Register() {
	var wg sync.WaitGroup
	for i := 1; i <= userNumber; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			register(i)
			token, err := login(strconv.Itoa(i))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			time.Sleep(10 * time.Millisecond)
			u := url.URL{Scheme: "ws", Host: addr, Path: "/api/connect/websocket"}
			err = Websocket(u.String(), token)
			if err != nil {
				fmt.Println(err.Error())
			}
		}(i)
		time.Sleep(10 * time.Millisecond)
	}
	wg.Wait()
}

type ReqLogin struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

type RspLogin struct {
	Data struct {
		UserInfo struct {
			UserId   string `json:"user_id"`
			Username string `json:"username"`
		} `json:"user_info"`
		Token string `json:"token"`
	}
}

func login(userId string) (token string, err error) {
	var req = ReqLogin{
		LoginName: "test-" + userId,
		Password:  "123456",
	}
	data, err := json.Marshal(req)
	if err != nil {
		return
	}
	res, err := Post(fmt.Sprintf("http://%v/api/user/login", addr), string(data))
	if err != nil {
		return
	}
	var rsp RspLogin
	if err = json.Unmarshal(res, &rsp); err != nil {
		return
	}
	fmt.Println(req.LoginName)
	token = rsp.Data.Token
	tokens = append(tokens, rsp.Data.Token)
	return
}

func Login() {
	var wg sync.WaitGroup
	for i := 1; i <= userNumber; i++ {
		wg.Add(1)
		go func(i string) {
			defer wg.Done()
			token, err := login(i)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			time.Sleep(10 * time.Millisecond)
			u := url.URL{Scheme: "ws", Host: addr, Path: "/api/connect/websocket"}
			err = Websocket(u.String(), token)
			if err != nil {
				fmt.Println(err.Error())
			}
		}(strconv.Itoa(i))
		time.Sleep(10 * time.Millisecond)
	}
	wg.Wait()
}

func detail(token string) {
	res, err := Get(fmt.Sprintf("http://%v/api/account/user/detail", addr), token, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(res))
}

func Detail() {
	var wg sync.WaitGroup
	for _, token := range tokens {
		wg.Add(1)
		go func(token string) {
			wg.Done()
			detail(token)
		}(token)
		time.Sleep(10 * time.Millisecond)
	}
	wg.Wait()
}

func WebSocket() {
	var wg sync.WaitGroup
	for _, token := range tokens {
		wg.Add(1)
		go func(token string) {
			defer wg.Done()
			u := url.URL{Scheme: "ws", Host: addr, Path: "/api/connect/websocket"}
			err := Websocket(u.String(), token)
			if err != nil {
				fmt.Println(err.Error())
			}
		}(token)
		time.Sleep(20 * time.Millisecond)
	}
	wg.Wait()
}

func main() {
	//Register()
	Login()
	//Detail()
	//WebSocket()
}
