package redis

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	str, err := Get("token:123")
	if err != nil {
		return
	}
	fmt.Println(str)
}

func TestAdd(t *testing.T)  {

}
