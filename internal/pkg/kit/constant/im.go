package constant

import "baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"

type OperatorDeleteStatus string

const (
	UnDelete      = "00"
	SendDelete    = "10"
	ReceiveDelete = "01"
	AllDelete     = "11"
)

const (
	JoinSessionOptReqContent = "已将加入会话请求发送给会话创建者，请耐心等候"
)

// 会话用户数限制
type SessionUserLimit int64

const (
	SessionBaseLevelUserLimit        = 2
	SessionNormalLevelUserLimit      = 20
	SessionSmallGroupLevelUserLimit  = 100
	SessionMediumGroupLevelUserLimit = 300
	SessionLargeGroupLevelUserLimit  = 500
	SessionTotalGroupLevelUserLimit  = 1000000
)

var SessionLevelUserLimitMap = map[im.SessionLevel]SessionUserLimit{
	im.SessionLevel_SessionBaseLevel:        SessionBaseLevelUserLimit,
	im.SessionLevel_SessionNormalLevel:      SessionNormalLevelUserLimit,
	im.SessionLevel_SessionSmallGroupLevel:  SessionSmallGroupLevelUserLimit,
	im.SessionLevel_SessionMediumGroupLevel: SessionMediumGroupLevelUserLimit,
	im.SessionLevel_SessionLargeGroupLevel:  SessionLargeGroupLevelUserLimit,
	im.SessionLevel_SessionTotalGroupLevel:  SessionTotalGroupLevelUserLimit,
}
