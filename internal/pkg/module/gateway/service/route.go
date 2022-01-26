package service

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	kitMiddleware "baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/module/gateway/config"
	"baby-fried-rice/internal/pkg/module/gateway/log"
	"baby-fried-rice/internal/pkg/module/gateway/middleware"
	"baby-fried-rice/internal/pkg/module/gateway/server"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func Register(engine *gin.Engine) {
	api := engine.Group("/api")

	//api.Use(middleware.ResponseHandle)
	api.Use(kitMiddleware.PreRequestHandle)
	api.POST("/admin/login", HandleManageProxy)
	api.POST("/user/register", HandleAccountUserProxy)
	api.POST("/user/login", HandleAccountUserProxy)

	user := api.Group("")
	user.Use(middleware.CheckToken)
	//user.Use(middleware.Auth)

	user.Any("/manage/*any", HandleManageProxy)
	user.Any("/account/user/*any", HandleAccountUserProxy)
	user.Any("/im/*any", HandleImProxy)
	user.Any("/space/*any", HandleSpaceProxy)
	user.Any("/comment/*any", HandleSpaceProxy)
	user.Any("/connect/*any", HandleConnectProxy)
	user.Any("/file/*any", HandleFileProxy)
	user.Any("/shop/*any", HandleShopProxy)
	user.Any("/live/*any", HandleLiveProxy)
	user.Any("/blog/*any", HandleBlogProxy)
}

func HandleAccountUserProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Rpc.SubServers.UserAccountServer)
}

func HandleManageProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Rpc.SubServers.ManageServer)
}

func HandleImProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Rpc.SubServers.ImServer)
}

func HandleSpaceProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Rpc.SubServers.SpaceServer)
}

func HandleConnectProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Rpc.SubServers.ConnectServer)
}

func HandleFileProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Rpc.SubServers.FileServer)
}

func HandleShopProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Rpc.SubServers.ShopServer)
}

func HandleLiveProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Rpc.SubServers.LiveServer)
}

func HandleBlogProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Rpc.SubServers.BlogServer)
}

func getLang(c context.Context) string {
	lang, ok := c.Value(handle.HeaderLanguage).(string)
	if !ok {
		lang = language.Chinese.String()
	}
	return lang
}

func handleProxy(c *gin.Context, serverName string) {
	serverUrl, err := server.GetRegisterClient().GetServer(serverName)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var parseServerUrl *url.URL
	if parseServerUrl, err = url.Parse(serverUrl); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parseServerUrl)
	if serverName != config.GetConfig().Rpc.SubServers.ConnectServer {
		proxy.ModifyResponse = func(response *http.Response) error {
			// 修改代理的返回值，如果错误码非0，则根据不同的语言添加不同的错误描述
			var body []byte
			body, err = ioutil.ReadAll(response.Body)
			if err != nil {
				log.Logger.Error(err.Error())
				return err
			}
			_ = response.Body.Close()

			var resp rsp.CommonResp
			if err = json.Unmarshal(body, &resp); err != nil {
				log.Logger.Error(err.Error())
				return err
			}
			if resp.Code != 0 {
				switch getLang(c) {
				case language.Chinese.String():
					resp.Message = constant.ErrCodeM[resp.Code]
				case language.English.String():
				default:
					resp.Message = constant.ErrCodeM[resp.Code]
				}
			}
			var data []byte
			data, err = json.Marshal(resp)
			if err != nil {
				log.Logger.Error(err.Error())
				return err
			}
			response.Body = ioutil.NopCloser(bytes.NewReader(data))
			response.ContentLength = int64(len(data))
			response.Header.Set("Content-Length", strconv.Itoa(len(data)))
			return nil
		}
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
