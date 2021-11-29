package handle

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func FileUpload(url, filePath string, userMeta UserMeta) (data []byte, err error) {
	var payload = &bytes.Buffer{}
	var writer = multipart.NewWriter(payload)
	var file *os.File
	if file, err = os.Open(filePath); err != nil {
		return
	}
	defer file.Close()
	var part io.Writer
	if part, err = writer.CreateFormFile("file", filepath.Base(filePath)); err != nil {
		return
	}
	if _, err = io.Copy(part, file); err != nil {
		return
	}
	if err = writer.Close(); err != nil {
		return
	}
	var client = &http.Client{}
	var req *http.Request
	if req, err = http.NewRequest("POST", url, payload); err != nil {
		return
	}
	req.Header.Set(HeaderAccountId, userMeta.AccountId)
	req.Header.Set(HeaderUsername, userMeta.Username)
	req.Header.Set(HeaderSchoolId, userMeta.SchoolId)
	req.Header.Set(HeaderReqId, userMeta.ReqId)
	req.Header.Set(HeaderIsOfficial, fmt.Sprintf("%v", userMeta.IsOfficial))
	req.Header.Set(HeaderPlatform, userMeta.Platform)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return
	}
	defer res.Body.Close()
	if data, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}
	return
}

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
