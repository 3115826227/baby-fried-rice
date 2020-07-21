package handle

import (
	"fmt"
	"testing"
)

func TestUpdateIp(t *testing.T) {
	describe := UpdateIp("112.10.111.212")
	fmt.Println(describe)
}

func TestGenerateSerialNumber(t *testing.T) {
	fmt.Println(GenerateSerialNumber())
}
