package handle

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func Post(url string, payload []byte, header http.Header) (data []byte, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		err = errors.Wrap(err, "failed to new request")
		return
	}
	req.Header = header
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "failed to client do")
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read all resp.body")
	}
	return
}

func Get(url string) (data []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to new request")
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read all resp.body")
		return
	}
	return
}
