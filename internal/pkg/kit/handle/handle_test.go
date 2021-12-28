package handle

import (
	"fmt"
	"testing"
)

func TestGeneratePhoneCode(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(GeneratePhoneCode())
	}
}
