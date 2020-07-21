package handle

import (
	"encoding/json"
	"github.com/3115826227/baby-fried-rice/module/im/src/config"
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func OfficialGroupAdd(c *gin.Context) {
	userMeta := GetUserMeta(c)
	var req model.ReqOfficialGroupAdd
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	data, err := Get(config.Config.AccountDaoUrl + "/dao/account/admin/organize/exist?id=" + req.Organize)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	var resp model.RspOkResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	if resp.Code != 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	var officialGroup = model.FriendGroupMeta{
		Name:               req.Name,
		Official:           true,
		SchoolDepartmentId: req.Organize,
		Origin:             userMeta.UserId,
	}
	id := GenerateID()
	now := time.Now()
	officialGroup.ID = id
	officialGroup.CreatedAt = now
	officialGroup.UpdatedAt = now

	var beans = make([]interface{}, 0)
	beans = append(beans, &officialGroup)
	beans = append(beans, &model.FriendGroupRelation{
		GroupId:         id,
		UserId:          userMeta.UserId,
		UserGroupRemark: userMeta.Username,
	})

	if err := db.CreateMulti(beans...); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	c.JSON(http.StatusOK, model.RspOkResponse{})
}

func GroupAdd(c *gin.Context) {
	var req model.GroupAddReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}
	userMeta := GetUserMeta(c)

	var group = model.FriendGroupMeta{
		Name:               req.Name,
		Level:              "0",
		Official:           false,
		SchoolDepartmentId: "",
		Origin:             userMeta.UserId,
	}
	id := GenerateID()
	now := time.Now()
	group.ID = id
	group.CreatedAt = now
	group.UpdatedAt = now

	var beans = make([]interface{}, 0)
	beans = append(beans, &group)
	for _, friend := range req.Friends {
		var relation = model.FriendGroupRelation{
			GroupId:         id,
			UserId:          friend.Id,
			UserGroupRemark: friend.Remark,
		}
		beans = append(beans, &relation)
	}
	beans = append(beans, &model.FriendGroupRelation{
		GroupId:         id,
		UserId:          userMeta.UserId,
		UserGroupRemark: userMeta.Username,
	})

	if err := db.CreateMulti(beans...); err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	c.JSON(http.StatusOK, model.RspOkResponse{})
}

func GroupListGet(c *gin.Context) {
	userMeta := GetUserMeta(c)

	var relations = make([]model.FriendGroupRelation, 0)
	if err := db.GetDB().Debug().Where("user_id = ?", userMeta.UserId).Find(&relations).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	var groupIds = make([]string, 0)
	for _, r := range relations {
		groupIds = append(groupIds, r.GroupId)
	}
	var groups = make([]model.FriendGroupMeta, 0)
	if err := db.GetDB().Debug().Where("id in (?)", groupIds).Find(&groups).Error; err != nil {
		log.Logger.Warn(err.Error())
	}

	SuccessResp(c, "", groups)
}
