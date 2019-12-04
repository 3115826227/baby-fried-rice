package handle

import (
	"fmt"
	"github.com/satori/go.uuid"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	u, _ := uuid.NewV4()
	fmt.Println(u.String())
	token, err := GenerateToken(u.String(), time.Now())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	userToken := ExplainToken(token)
	fmt.Println(userToken)
	fmt.Println(token)
}

func TestExplainToken(t *testing.T) {
	userToken := ExplainToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVfdGltZSI6IjIwMTktMTEtMTZUMTk6Mjk6MjUuNzkxMDc0KzA4OjAwIiwidXNlcl9pZCI6ImU0NDE1ODlhLTgxMzItNGUxMy04Yzk5LWUwYTk4ODBjMWQ5ZiJ9.TYYrf-SQcvTsskG_VspuRrBpGxwI1TAxVFSjQYmv2Ms")
	fmt.Println(userToken)
}

func TestEncodePassword(t *testing.T) {
	encodePassword := EncodePassword("root")
	fmt.Println(encodePassword)
}

func TestAddRoot(t *testing.T) {
	AddRoot()
}
