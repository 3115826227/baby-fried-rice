package redis

import "testing"

func TestAddAccountToken(t *testing.T) {
	AddAccountToken("token:123", "hello")
}