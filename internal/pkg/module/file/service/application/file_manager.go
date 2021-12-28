package application

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/file/log"
)

type FileManager interface {
	UploadFile(key, localFilePath string) (ossMeta tables.OssMeta, downUrl string, err error)
	DeleteFile(bucket, key string) error
}

var (
	fileManager FileManager
)

func InitFileManager(fileMode constant.FileMode) (err error) {
	switch fileMode {
	case constant.LocalFileMode:
		fileManager, err = NewMinIOManager(log.Logger, 1)
	case constant.QiNiuYunOssFileMode:
		fileManager, err = NewOssManager(log.Logger)
	}
	return
}

func GetFileManager() FileManager {
	return fileManager
}
