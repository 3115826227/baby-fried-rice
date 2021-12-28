package application

import (
	"baby-fried-rice/internal/pkg/module/file/log"
	"testing"
)

func TestNewMinIOManager(t *testing.T) {
	err := log.InitLog("file", "DEBUG", "")
	if err != nil {
		panic(err)
	}
	manager, err := NewMinIOManager(log.Logger, 1)
	if err != nil {
		panic(err)
	}
	_, _, err = manager.UploadFile("Dockerfile", "Dockerfile")
	if err != nil {
		panic(err)
	}

}
