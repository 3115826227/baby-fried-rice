package model

import (
	"fmt"
	"testing"
)

func TestGetFriend(t *testing.T) {
	res, err := GetFriend("7e2db36e-aff3-4793-b65a-bc2580a9ccec")
	if err != nil {
		return
	}
	fmt.Println(res)
}
