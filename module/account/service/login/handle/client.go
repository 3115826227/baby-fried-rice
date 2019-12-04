package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/account/service/model"
	"net/http"
	"time"
	"github.com/3115826227/baby-fried-rice/module/account/service/model/db"
	"github.com/jinzhu/gorm"
	"github.com/3115826227/baby-fried-rice/module/account/log"
	"fmt"
)

func IsDuplicateClient(c *gin.Context, name string, isAdd bool) bool {
	var client = new(model.AccountClient)
	err := db.DB.First(&client, map[string]interface{}{"name": name}).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return true
	}
	layout := ""
	if isAdd {
		layout = `"%s"已存在,无法新增客户,请重新输入`
	} else {
		layout = `"%s"已存在,无法编辑客户,请重新输入`
	}
	if client.ID != "" {
		ErrorResp(c, http.StatusBadRequest, ErrCodeDuplicateName,
			fmt.Sprintf(layout, name))
		return true
	}
	return false
}

func ClientAdd(c *gin.Context) {
	var err error
	var req = model.ReqClientAdd{}
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	if IsDuplicateClient(c, req.Name, true) {
		return
	}

	var schools = make([]model.School, 0)
	err = db.DB.Find(&schools).Where("id in (?)", req.SchoolIds).Error
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	if len(schools) < len(req.SchoolIds) {
		ErrorResp(c, http.StatusBadRequest, ErrSchoolIdNotFound, ErrCodeM[ErrSchoolIdNotFound])
		return
	}

	var beans = make([]interface{}, 0)

	var userMeta = GetUserMeta(c)
	var client = new(model.AccountClient)
	var now = time.Now()
	client.ID = GenerateID()
	client.Name = req.Name
	client.CreatedAt, client.UpdatedAt = now, now
	client.Origin = userMeta.UserId

	for _, id := range req.SchoolIds {
		var rel = model.ClientSchoolRelation{
			SchoolId: id,
			ClientId: client.ID,
		}
		beans = append(beans, &rel)
	}
	beans = append(beans, &client)

	if err = db.CreateMulti(beans...); err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}

func ClientDelete(c *gin.Context) {

}
