package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/game"
	"baby-fried-rice/internal/pkg/module/game/grpc"
	"baby-fried-rice/internal/pkg/module/game/log"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 个人的游戏对局列表查询
func GameRecordQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var gameType int
	if gameType, err = strconv.Atoi(c.Query("game_type")); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var gameClient game.DaoGameClient
	if gameClient, err = grpc.GetGameClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var reqGame = game.ReqGameRecordQueryDao{
		GameType:  game.GameType(gameType),
		AccountId: userMeta.AccountId,
		Page:      reqPage.Page,
		PageSize:  reqPage.PageSize,
	}
	var resp *game.RspGameRecordQueryDao
	resp, err = gameClient.GameRecordQueryDao(context.Background(), &reqGame)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	for _, item := range resp.List {
		var record = rsp.GameRecord{
			GameRecordId:    item.GameRecordId,
			GameStatus:      item.GameStatus,
			GameType:        item.GameType,
			Result:          item.GameResult,
			FinishTimestamp: item.FinishTimestamp,
			UserRole:        item.UserRole,
		}
		list = append(list, record)
	}
	handle.SuccessListResp(c, "", list, resp.Total, resp.Page, resp.PageSize)
}

// 中国象棋游戏详情数据查询
func GameRecordDetailQueryHandle(c *gin.Context) {

}

// 中国象棋游戏状态数据查询
func ChinaChessGameStatusDataQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	gameRecordId, err := strconv.Atoi(c.Query("game_record_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var gameType int
	if gameType, err = strconv.Atoi(c.Query("game_type")); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var gameClient game.DaoGameClient
	if gameClient, err = grpc.GetGameClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var reqGame = game.ReqGameStatusQueryDao{
		GameRecordId: int64(gameRecordId),
		GameType:     game.GameType(gameType),
		AccountId:    userMeta.AccountId,
	}
	var resp *game.RspGameStatusQueryDao
	if resp, err = gameClient.GameStatusQueryDao(context.Background(), &reqGame); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var board rsp.ChinaChessBoard
	if err = json.Unmarshal([]byte(resp.GameStatusData), &board); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var response = rsp.ChinaChessStatusResp{
		UserRole:   resp.UserRole,
		GameStatus: resp.GameStatus,
		Board:      board,
	}
	handle.SuccessResp(c, "", response)
}
