package middleware

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	var resp rsp.CommonResp
	_ = json.Unmarshal(b, &resp)
	if resp.Code != 0 {
		resp.Message = constant.ErrCodeM[resp.Code]
	}
	data, _ := json.Marshal(resp)
	w.ResponseWriter.Header().Set("Content-Length", fmt.Sprintf("%v", len(data)))
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// 返回值处理 todo 会出现panic 暂未解决
func ResponseHandle(c *gin.Context) {
	w := &responseBodyWriter{
		body:           &bytes.Buffer{},
		ResponseWriter: c.Writer,
	}
	c.Writer = w
	c.Next()
}
