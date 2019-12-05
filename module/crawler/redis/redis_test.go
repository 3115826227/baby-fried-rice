package redis

import (
	"fmt"
	"testing"
)

func TestExist(t *testing.T) {
	Add("hello", "world")
	fmt.Println(Exist("hello"))
}
