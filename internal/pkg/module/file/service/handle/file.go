package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/module/file/db"
	"baby-fried-rice/internal/pkg/module/file/log"
	"baby-fried-rice/internal/pkg/module/file/service/application"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	DefaultFileStorageDay = 3
)

func FileUploadHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}

	var data []byte
	data, err = ioutil.ReadAll(file)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	defer file.Close()

	if err = ioutil.WriteFile(header.Filename, data, 0755); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	defer func() {
		if err = os.Remove(fmt.Sprintf("./%v", header.Filename)); err != nil {
			log.Logger.Error(err.Error())
		}
	}()

	now := time.Now()
	var localPath = header.Filename
	var ossMeta tables.OssMeta
	var downUrl string
	ossMeta, downUrl, err = application.GetFileManager().UploadFile(header.Filename, localPath)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}

	var f = tables.File{
		Origin:         userMeta.AccountId,
		PermissionType: 1,
		FileName:       header.Filename,
		FileType:       1,
		FileSize:       header.Size,
		DownUrl:        downUrl,
		StorageDay:     DefaultFileStorageDay,
		Bucket:         ossMeta.Bucket,
	}
	f.ID = handle.GenerateSerialNumberByLen(10)
	f.CreatedAt = now
	if err = db.GetDB().CreateObject(&f); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}

	var resp = rsp.FileUploadResp{
		File: rsp.File{
			ID:         f.ID,
			Origin:     userMeta.AccountId,
			Name:       header.Filename,
			DownUrl:    downUrl,
			Size:       header.Size,
			UploadTime: now.Unix(),
			StorageDay: DefaultFileStorageDay,
		},
	}
	handle.SuccessResp(c, "", resp)
}

func FileQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var files = make([]tables.File, 0)
	if err := db.GetDB().GetDB().Where("origin = ?", userMeta.AccountId).Find(&files).Error; err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var list = make([]rsp.File, 0)
	for _, f := range files {
		var file = rsp.File{
			ID:         f.ID,
			Origin:     f.Origin,
			Name:       f.FileName,
			DownUrl:    f.DownUrl,
			Size:       f.FileSize,
			UploadTime: f.UpdatedAt.Unix(),
			StorageDay: f.StorageDay,
		}
		list = append(list, file)
	}
	var resp = rsp.FileQueryResp{List: list}
	handle.SuccessResp(c, "", resp)
}

func FileDeleteHandle(c *gin.Context) {
	id := c.Query("id")
	userMeta := handle.GetUserMeta(c)
	var file tables.File
	if err := db.GetDB().GetDB().Where("id = ? and origin = ?",
		id, userMeta.AccountId).First(&file).Error; err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	tx := db.GetDB().GetDB().Begin()
	if err := tx.Where("id = ? and origin = ?", id, userMeta.AccountId).Delete(&tables.File{}).Error; err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	if err := application.GetFileManager().DeleteFile(file.Bucket, file.FileName); err != nil {
		tx.Rollback()
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	tx.Commit()
	handle.SuccessResp(c, "", nil)
}
