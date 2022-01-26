package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/space/grpc"
	"baby-fried-rice/internal/pkg/module/space/log"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func sendAddSpaceNotify(space rsp.SpaceResp, accountId string) {
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var detailResp *user.RspDaoUserDetail
	detailResp, err = userClient.UserDaoDetail(context.Background(), &user.ReqDaoUserDetail{AccountId: accountId})
	var userResp *user.RspUserDaoAll
	userResp, err = userClient.UserDaoAll(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var now = time.Now().Unix()
	for _, id := range userResp.AccountIds {
		var notify = models.WSMessageNotify{
			WSMessageNotifyType: constant.SpaceMessageNotify,
			Receive:             id,
			WSMessage: models.WSMessage{
				Space: &rsp.SpaceResp{
					Id:          space.Id,
					Content:     space.Content,
					Images:      space.Images,
					VisitorType: space.VisitorType,
					Origin: &rsp.User{
						AccountID:  detailResp.Detail.AccountId,
						Username:   detailResp.Detail.Username,
						HeadImgUrl: detailResp.Detail.HeadImgUrl,
						IsOfficial: detailResp.Detail.IsOfficial,
					},
					CreateTime: now,
				},
				Send: rsp.User{
					AccountID:  detailResp.Detail.AccountId,
					Username:   detailResp.Detail.Username,
					HeadImgUrl: detailResp.Detail.HeadImgUrl,
					IsOfficial: detailResp.Detail.IsOfficial,
				},
			},
			Timestamp: now,
		}
		if err = mq.Send(topic, notify.ToString()); err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}
}

func sendCommentNotify() {

}

func sendOptNotify() {

}
