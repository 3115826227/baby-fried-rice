package handle

import "github.com/gin-gonic/gin"

type ShopRegisterReq struct {
	LoginName string `json:"login_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Describe  string
}

type ShopUpdateReq struct {
	ShopId   int    `json:"shop_id" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Describe string `json:"describe" binding:"required"`
}

//商家注册
func ShopRegisterHandle(c *gin.Context) {

}

//商家信息更新
func ShopUpdateHandle(c *gin.Context) {

}
