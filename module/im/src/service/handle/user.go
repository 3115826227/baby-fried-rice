package handle

import "github.com/gin-gonic/gin"

type AddUserReq struct {
	UserId    string `json:"user_id"`
	AccountId string `json:"account_id"`
}

func AddUser(c *gin.Context) {

}
