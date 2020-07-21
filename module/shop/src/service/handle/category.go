package handle

import "github.com/gin-gonic/gin"

type CommodityCategoryAddReq struct {
	ParentId string `json:"parent_id"`
	Category string `json:"category" binding:"required"`
}

type CommodityCategoryUpdateReq struct {
	CategoryId int    `json:"category_id" binding:"required"`
	Category   string `json:"category" binding:"required"`
}

func CommodityCategoryAddHandle(c *gin.Context) {

}

func CommodityCategoryUpdateHandle(c *gin.Context) {

}
