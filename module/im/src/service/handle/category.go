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
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	var category model.FriendCategoryMeta
	userMeta := GetUserMeta(c)
	category.Origin = userMeta.UserId
	category.Name = req.Name
	if err := db.GetDB().Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}

/*
	修改好友分类
*/
func FriendCategoryUpdate(c *gin.Context) {
	var req model.FriendCategoryUpdateReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	_, err := model.FindFriendCategory(req.CategoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	var updateMap = map[string]interface{}{"name": req.Name}
	if err := db.GetDB().Model(model.FriendCategoryMeta{}).Where("id = ?", req.CategoryId).Updates(updateMap).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
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
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}
	category, err := model.FindFriendCategory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}
	relations, err := model.GetFriendCategoryRelation(category.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}
	var beans = make([][]interface{}, 0)
	beans = append(beans, []interface{}{category})
	for _, relation := range relations {
		beans = append(beans, []interface{}{relation})
	}
	if err := db.DeleteMulti(beans); err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}
