package register

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReqRegister struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Tag     string `json:"tag"`
	Name    string `json:"name"`
}

func Register(c *gin.Context) {
	var req ReqRegister
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

}
