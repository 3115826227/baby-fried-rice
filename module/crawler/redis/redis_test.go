package redis

import (
	"testing"
	"fmt"
)

func TestExist(t *testing.T) {
	Add("hello", "world")
	fmt.Println(Exist("hello"))
}
