package handle

import (
	"fmt"
	"testing"
)

func TestAddGrade(t *testing.T) {
	AddGrade()
}

func TestAddSubject(t *testing.T) {
	AddSubject()
}

func TestAddCourse(t *testing.T) {
	AddCourse()
}

func TestIsValidCourse(t *testing.T) {
	fmt.Println(IsValidCourse(10, 1))
}

func TestIsValidArea(t *testing.T) {
	fmt.Println(IsValidArea("330110"))
}
