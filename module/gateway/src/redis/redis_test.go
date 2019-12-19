package redis

import (
	"testing"
	"fmt"
)

func TestGet(t *testing.T) {
	str, err := Get("token:123")
	if err != nil {
		return
	}
	fmt.Println(str)
}
