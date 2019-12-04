package utils

import (
	"net/http"
	"github.com/3115826227/baby-fried-rice/module/public/log"
	"io/ioutil"
	"strings"
)

func Request(url string) (data []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	return
}

func PostRequest(method, url string, body *strings.Reader) (data []byte, err error) {
	c := http.Client{}
	var req *http.Request
	if body == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, body)
	}
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	return
}
