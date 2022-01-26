package middleware

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func PreRequestHandle(c *gin.Context) {
	// 生成uuid
	u := uuid.NewV4().String()
	c.Request.Header.Set(handle.HeaderUUID, u)
	c.Set(handle.HeaderUUID, u)
	// 语言解析
	lang := c.GetHeader(handle.HeaderLanguage)
	c.Request = c.Request.WithContext(context.WithValue(c, handle.HeaderLanguage, lang))
	c.Set(handle.HeaderLanguage, lang)
	c.Next()
}
