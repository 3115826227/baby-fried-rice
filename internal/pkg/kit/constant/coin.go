package constant

const (
	MaxTimestamp       = 1e12
	UserCoinBoardTotal = 20
	UserRegisterCoin   = 500
)

// 积分使用类型
type CoinType int64

const (
	// 签到
	SignInCoinType CoinType = 1
	// 消费
	ConsumeCoinType = 101
	// 新用户注册赠送
	UserRegisterCoinType = 201
	// 系统赠送
	SystemGiveawayCoinType = 301
)

const (
	SignInCoinDescribe         = "签到成功获取积分"
	ConsumeCoinDescribe        = "消费使用消耗积分"
	UserRegisterCoinDescribe   = "用户注册成功获取积分"
	SystemGiveawayCoinDescribe = "系统赠送获取积分"
)

var CoinTypeDescribeMap = map[CoinType]string{
	SignInCoinType:         SignInCoinDescribe,
	ConsumeCoinType:        ConsumeCoinDescribe,
	UserRegisterCoinType:   UserRegisterCoinDescribe,
	SystemGiveawayCoinType: SystemGiveawayCoinDescribe,
}
