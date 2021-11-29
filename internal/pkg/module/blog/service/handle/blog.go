package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/blog"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/blog/config"
	"baby-fried-rice/internal/pkg/module/blog/grpc"
	"baby-fried-rice/internal/pkg/module/blog/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 添加标签
func TagAddHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.ReqAddTag
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if _, err = client.TagAddDao(context.Background(),
		&blog.ReqTagAddDao{Origin: userMeta.AccountId, Tags: []string{req.Tag}}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 查询标签
func TagHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *blog.RspTagQueryDao
	if resp, err = client.TagQueryDao(context.Background(),
		&blog.ReqTagQueryDao{Origin: userMeta.AccountId}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var response = rsp.TagResp{Tags: resp.Tags}
	handle.SuccessResp(c, "", response)
}

// 删除标签
func TagDeleteHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	tags := strings.Split(c.Query("tags"), ",")
	var err error
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if _, err = client.TagDeleteDao(context.Background(),
		&blog.ReqTagDeleteDao{Origin: userMeta.AccountId, Tags: tags}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 添加分类
func CategoryAddHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.ReqAddCategory
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if _, err = client.CategoryAddDao(context.Background(),
		&blog.ReqCategoryAddDao{Origin: userMeta.AccountId, Category: req.Category}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 查询分类
func CategoryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *blog.RspCategoryQueryDao
	if resp, err = client.CategoryQueryDao(context.Background(),
		&blog.ReqCategoryQueryDao{Origin: userMeta.AccountId}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var response = rsp.CategoryResp{Categories: resp.Category}
	handle.SuccessResp(c, "", response)
}

// 删除分类
func CategoryDeleteHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	category := c.Query("category")
	var err error
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if _, err = client.CategoryDeleteDao(context.Background(),
		&blog.ReqCategoryDeleteDao{Origin: userMeta.AccountId, Category: category}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func writeFile(data []byte, fileName string) error {
	return ioutil.WriteFile(fileName, data, 755)
}

func deleteFile(fileName string) error {
	return os.Remove(fileName)
}

// 添加博文
func BlogAddHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.ReqAddBlog
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var blogReq = blog.ReqBlogAddDao{
		Origin:         userMeta.AccountId,
		Title:          req.Title,
		Tags:           req.Tags,
		Category:       req.Category,
		PreviewContent: req.Content,
	}
	if _, err = client.BlogAddDao(context.Background(), &blogReq); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var fileName = fmt.Sprintf("file/blog_%v_%v.md", userMeta.AccountId, req.Title)
	if err = writeFile([]byte(req.Content), fileName); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var data []byte
	if data, err = handle.FileUpload(config.GetConfig().Rpc.SubServers.FileServer, fileName, *userMeta); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	log.Logger.Info(string(data))
	defer func() {
		if err = deleteFile(fileName); err != nil {
			log.Logger.Error(err.Error())
		}
	}()
	handle.SuccessResp(c, "", nil)
}

// 更新博文
func BlogUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.ReqUpdateBlog
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var blogReq = blog.ReqBlogUpdateDao{
		Origin:         userMeta.AccountId,
		Title:          req.Title,
		Tags:           req.Tags,
		Category:       req.Category,
		PreviewContent: req.Content,
		Status:         req.Status,
	}
	if _, err = client.BlogUpdateDao(context.Background(), &blogReq); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var fileName = fmt.Sprintf("file/blog_%v_%v.md", userMeta.AccountId, req.Title)
	if err = writeFile([]byte(req.Content), fileName); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var data []byte
	if data, err = handle.FileUpload(config.GetConfig().Rpc.SubServers.FileServer, fileName, *userMeta); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	log.Logger.Info(string(data))
	defer func() {
		if err = deleteFile(fileName); err != nil {
			log.Logger.Error(err.Error())
		}
	}()
	handle.SuccessResp(c, "", nil)
}

// 查询博文列表
func BlogHandle(c *gin.Context) {
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var queryType int
	if queryType, err = strconv.Atoi(c.Query("query_type")); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var blogReq = blog.ReqBlogQueryDao{
		Search:      c.Query("search"),
		BloggerLike: c.Query("blogger_like"),
		QueryType:   blog.QueryType(queryType),
		Page:        reqPage.Page,
		PageSize:    reqPage.PageSize,
	}
	var resp *blog.RspBlogQueryDao
	if resp, err = client.BlogQueryDao(context.Background(), &blogReq)
		err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var ids = make([]string, 0)
	for _, b := range resp.List {
		ids = append(ids, b.Blogger)
	}
	var userClient user.DaoUserClient
	if userClient, err = grpc.GetUserClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	if userResp, err = userClient.UserDaoById(context.Background(),
		&user.ReqUserDaoById{Ids: ids}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userMap = make(map[string]rsp.User)
	for _, u := range userResp.Users {
		userMap[u.Id] = rsp.User{
			AccountID:  u.Id,
			Username:   u.Username,
			HeadImgUrl: u.HeadImgUrl,
			IsOfficial: u.IsOfficial,
		}
	}
	var list = make([]interface{}, 0)
	for _, b := range resp.List {
		list = append(list, rsp.Blog{
			BlogId:         b.BlogId,
			Blogger:        userMap[b.Blogger],
			Title:          b.Title,
			Category:       b.Category,
			PreviewContent: b.PreviewContent,
			LikeTotal:      b.LikeTotal,
			ReadTotal:      b.ReadTotal,
			CommentTotal:   b.CommentTotal,
			Timestamp:      b.Timestamp,
		})
	}
	handle.SuccessListResp(c, "", list, resp.Total, resp.Page, resp.PageSize)
}

// 查询博文详情
func BlogDetailHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	blogId := c.Query("blog_id")
	var err error
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *blog.RspBlogDetailQueryDao
	if resp, err = client.BlogDetailQueryDao(context.Background(),
		&blog.ReqBlogDetailQueryDao{BlogId: blogId, AccountId: userMeta.AccountId}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userClient user.DaoUserClient
	if userClient, err = grpc.GetUserClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	if userResp, err = userClient.UserDaoById(context.Background(),
		&user.ReqUserDaoById{Ids: []string{resp.Blog.Blogger}}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if len(userResp.Users) != 1 {
		err = fmt.Errorf("failed to get blogger by user dao")
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var u = userResp.Users[0]
	var response = rsp.BlogDetailResp{
		Blog: rsp.Blog{
			BlogId: resp.Blog.BlogId,
			Blogger: rsp.User{
				AccountID:  u.Id,
				Username:   u.Username,
				HeadImgUrl: u.HeadImgUrl,
				IsOfficial: u.IsOfficial,
			},
			Title:          resp.Blog.Title,
			Category:       resp.Blog.Category,
			PreviewContent: resp.Blog.PreviewContent,
			LikeTotal:      resp.Blog.LikeTotal,
			ReadTotal:      resp.Blog.ReadTotal,
			CommentTotal:   resp.Blog.CommentTotal,
			Timestamp:      resp.Blog.Timestamp,
		},
	}
	handle.SuccessResp(c, "", response)
}

// 删除博文
func BlogDeleteHandle(c *gin.Context) {
	ids := strings.Split(c.Query("ids"), ",")
	userMeta := handle.GetUserMeta(c)
	var err error
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if _, err = client.BlogDeleteDao(context.Background(),
		&blog.ReqBlogDeleteDao{Origin: userMeta.AccountId, BlogIds: ids}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 查询博主
func BloggerHandle(c *gin.Context) {
	blogger := c.Query("blogger")
	var err error
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *blog.RspBloggerQueryDao
	if resp, err = client.BloggerQueryDao(context.Background(),
		&blog.ReqBloggerQueryDao{Blogger: blogger}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userClient user.DaoUserClient
	if userClient, err = grpc.GetUserClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	if userResp, err = userClient.UserDaoById(context.Background(),
		&user.ReqUserDaoById{Ids: []string{resp.Blogger}}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if len(userResp.Users) != 1 {
		err = fmt.Errorf("failed to get blogger by user dao")
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var u = userResp.Users[0]
	var response = rsp.BloggerResp{
		Blogger: rsp.User{
			AccountID:  u.Id,
			Username:   u.Username,
			HeadImgUrl: u.HeadImgUrl,
			IsOfficial: u.IsOfficial,
		},
		Tags:      resp.Tags,
		BlogTotal: resp.BlogTotal,
		LikeTotal: resp.LikeTotal,
		ReadTotal: resp.ReadTotal,
		FansTotal: resp.FansTotal,
	}
	handle.SuccessResp(c, "", response)
}

// 查询粉丝
func BloggerFansHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *blog.RspFansQueryDao
	if resp, err = client.FansQueryDao(context.Background(),
		&blog.ReqFansQueryDao{
			Blogger:  userMeta.AccountId,
			Page:     reqPage.Page,
			PageSize: reqPage.PageSize,
		}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userClient user.DaoUserClient
	if userClient, err = grpc.GetUserClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	if userResp, err = userClient.UserDaoById(context.Background(),
		&user.ReqUserDaoById{Ids: resp.Users}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userMap = make(map[string]rsp.User)
	for _, u := range userResp.Users {
		userMap[u.Id] = rsp.User{
			AccountID:  u.Id,
			Username:   u.Username,
			HeadImgUrl: u.HeadImgUrl,
			IsOfficial: u.IsOfficial,
		}
	}
	var list = make([]interface{}, 0)
	for _, u := range resp.Users {
		list = append(list, userMap[u])
	}
	handle.SuccessListResp(c, "", list, resp.Total, resp.Page, resp.PageSize)
}

// 关注/取关博主
func BloggerFocusOnHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.ReqFocusOnBlogger
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if _, err = client.FocusAddDao(context.Background(),
		&blog.ReqFocusAddDao{
			Origin:  req.Blogger,
			Focus:   req.FocusOn,
			Account: userMeta.AccountId,
		}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 喜欢/取消喜欢博主
func BlogLikeHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.ReqBlogLikeBlogger
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var client blog.DaoBlogClient
	if client, err = grpc.GetBlogClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if _, err = client.BlogLikeAddDao(context.Background(),
		&blog.ReqBlogLikeAddDao{
			BlogId:  req.BlogId,
			Like:    req.Like,
			Account: userMeta.AccountId,
		}); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}
