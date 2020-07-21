package handle

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
)

/*
	添加好友分类
*/
func FriendCategoryAdd(c *gin.Context) {
	var req model.FriendCategoryAddReq
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var category model.FriendCategoryMeta
	userMeta := GetUserMeta(c)
	category.Origin = userMeta.UserId
	category.Name = req.Name
	if err := db.GetDB().Debug().Create(&category).Error; err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}

/*
	修改好友分类名称
*/
func FriendCategoryUpdate(c *gin.Context) {
	var req model.FriendCategoryUpdateReq
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	if _, err := model.FindFriendCategory(req.CategoryId); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var updateMap = map[string]interface{}{"name": req.Name}
	if err := db.GetDB().Debug().Model(model.FriendCategoryMeta{}).Where("id = ?", req.CategoryId).Updates(updateMap).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}

/*
	获取好友分类列表
*/
func FriendCategory(c *gin.Context) {
	userMeta := GetUserMeta(c)
	category, err := model.GetFriendCategory(userMeta.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	sort.Sort(model.FriendCategories(category))
	SuccessResp(c, "", category)
}

/*
	删除好友分类
*/
func FriendCategoryDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	category, err := model.FindFriendCategory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	relations, err := model.GetFriendCategoryRelation(category.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	if len(relations) != 0 {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	if err := db.GetDB().Debug().Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}
