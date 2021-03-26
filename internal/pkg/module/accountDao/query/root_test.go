package query

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"fmt"
	"testing"
)

func TestGetRootByLogin(t *testing.T) {
	root, err := GetRootByLogin("root", handle.EncodePassword("root"))
	if err != nil {
		panic(err)
	}
	fmt.Println(root)
}
