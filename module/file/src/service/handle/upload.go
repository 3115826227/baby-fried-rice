package handle

import (
	"context"
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/file/src/log"
	"github.com/3115826227/baby-fried-rice/module/file/src/model"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadHandle(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "文件上传失败"})
		return
	}

	var data []byte
	data, err = ioutil.ReadAll(file)
	defer file.Close()
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "文件读取失败"})
		return
	}
	if err = ioutil.WriteFile(header.Filename, data, 0755); err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "文件保存失败"})
		return
	}
	defer func() {
		if err = os.Remove(fmt.Sprintf("./%v", header.Filename)); err != nil {
			log.Logger.Warn(err.Error())
		}
	}()

	var ossInfo model.OssInfo
	if ossInfo, err = DataUpOss(header.Size, fmt.Sprintf("./%v", header.Filename), header.Filename); err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "文件上传失败"})
		return
	}

	downUrl := fmt.Sprintf("%v/%v", ossInfo.Domain, header.Filename)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "上传成功",
		"data": gin.H{
			"name":        header.Filename,
			"down_url":    downUrl,
			"size":        header.Size,
			"upload_time": ossInfo.UpdatedAt,
			"storage_day": 3,
		},
	})
}

func DataUpOss(size int64, filePath, key string) (ossInfo model.OssInfo, err error) {
	ossInfo = model.OssInfoQuery()
	var putPolicy storage.PutPolicy
	putPolicy.Scope = ossInfo.Bucket

	mac := qbox.NewMac(ossInfo.AccessKey, ossInfo.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	cfg.UseHTTPS = false
	cfg.Zone = &storage.ZoneHuadong
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	if err = formUploader.PutFile(context.Background(), &ret, upToken, key, filePath, &putExtra); err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	go model.OssInfoDataUpdate(ossInfo, size)
	return
}
