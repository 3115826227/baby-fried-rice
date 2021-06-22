package constant

const (
	SessionUserCount = 200
)

type WSMessageNotifyType int

const (
	// 私信消息通知
	PrivateMessageNotify = 1
	// 会话消息通知
	SessionMessageNotify = 2
	// 动态消息通知
	SpaceMessageNotify = 3
)

type WSMessageType int

const (
	// 文本消息
	TextMessage = 1
	// 文件消息
	FileMessage = 2
	// 图片消息
	ImgMessage = 3
)
