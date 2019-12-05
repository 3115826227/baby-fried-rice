package utils

import (
	"fmt"
	"testing"
)

func TestRequest(t *testing.T) {
	data, err := Request("http://www.baidu.com")
	if err != nil {
		return
	}
	fmt.Println(string(data))
}
