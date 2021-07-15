package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/module/backend/cache"
	"baby-fried-rice/internal/pkg/module/backend/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

/*
	缓存监控信息
*/
func CacheInfoHandle(c *gin.Context) {
	result, err := cache.GetCache().Info()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	results := strings.Split(result, "\n")
	var mp = make(map[string]interface{})
	for _, res := range results {
		if strings.Contains(res, ":") {
			slice := strings.Split(res, ":")
			key := slice[0]
			slice[1] = strings.Trim(slice[1], "\r")
			var val interface{}
			val, err = strconv.ParseFloat(slice[1], 64)
			if err != nil {
				val = slice[1]
			}
			mp[key] = val
		}
	}
	handle.SuccessResp(c, "", mp)
}
