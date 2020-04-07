package handle

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	DefaultCategory = "默认分组"
)

/*
	添加好友
*/
func FriendAdd(c *gin.Context) {
	var req model.FriendAddReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}
	//todo 对userId的验证
	userMeta := GetUserMeta(c)

	var friend = model.FriendRelation{}
	var relativeFriend = model.FriendRelation{}
	friend.ID = GenerateID()
	friend.Origin = userMeta.UserId
	friend.Friend = req.UserId
	friend.FriendRemark = req.Remark

	relativeFriend.ID = GenerateID()
	relativeFriend.Origin = req.UserId
	relativeFriend.Friend = userMeta.UserId
	relativeFriend.FriendRemark = userMeta.Username

	var beans = make([]interface{}, 0)
	beans = append(beans, &friend)
	beans = append(beans, &relativeFriend)

	if err := db.CreateMulti(beans...); err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	c.JSON(http.StatusOK, model.RspOkResponse{})
}

/*
	修改好友备注
*/
func FriendRemarkUpdate(c *gin.Context) {
	var req model.FriendRemarkUpdateReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	_, err := model.FindFriendRelationById(req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}
	var updateMap = map[string]interface{}{"friend_remark": req.Remark}
	if err := db.GetDB().Model(model.FriendRelation{}).Where("id = ?", req.Id).Update(updateMap).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	c.JSON(http.StatusOK, model.RspOkResponse{})
}

/*
	获取好友列表
*/
func Friends(c *gin.Context) {
	userMeta := GetUserMeta(c)
	res, err := model.GetFriend(userMeta.UserId)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	var rsp = make([]model.RspFriendCategory, 0)
	var categoryMap = make(map[string][]model.RspFriend, 0)
	for _, friend := range res {
		if friend.CategoryName == "" {
			friend.CategoryName = DefaultCategory
		}
		if _, exist := categoryMap[friend.CategoryName]; !exist {
			categoryMap[friend.CategoryName] = make([]model.RspFriend, 0)
		}
		friends := categoryMap[friend.CategoryName]
		friends = append(friends, model.RspFriend{
			Id:     friend.Id,
			Friend: friend.Friend,
			Remark: friend.FriendRemark,
		})
		categoryMap[friend.CategoryName] = friends
	}
	for category, friends := range categoryMap {
		rsp = append(rsp, model.RspFriendCategory{
			Name:    category,
			Friends: friends,
		})
	}

	SuccessResp(c, "", rsp)
}

/*
	删除好友
*/
func FriendDelete(c *gin.Context) {
	id := c.Query("id")
	friend, err := model.FindFriendRelationById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}
	userMeta := GetUserMeta(c)
	relativeFriend, err := model.FindFriendRelation(friend.Friend, userMeta.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	var beans = make([][]interface{}, 0)
	beans = append(beans, []interface{}{&friend})
	beans = append(beans, []interface{}{&relativeFriend})

	relations, err := model.FindFriendCategoryRelationByRelationId(relativeFriend.ID, friend.ID)
	if err == nil {
		for _, relation := range relations {
			beans = append(beans, []interface{}{&relation})
		}
	}

	if err := db.DeleteMulti(beans); err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	c.JSON(http.StatusOK, model.RspOkResponse{})
}
