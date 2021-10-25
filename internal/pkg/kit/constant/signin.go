package constant

type SignInType int64

const (
	// 正常签到
	NormalSignInType SignInType = 0
	// 补签
	MakeupSignInType = 1
)

type RewardCoinBySignedInType int64

/*
	连续签到积分奖励
		连续签到一天： 奖励积分20
		连续签到两天： 奖励积分30
		连续签到三天： 奖励积分50
		连续签到四天： 奖励积分80
		连续签到五天： 奖励积分120
		连续签到六天： 奖励积分200
		连续签到七天： 奖励积分300
	连续签到七天之后，第八天按连续签到一天重新计算
*/
const (
	RewardCoinBySignedInContinuedOneDay   RewardCoinBySignedInType = 20
	RewardCoinBySignedInContinuedTwoDay                            = 30
	RewardCoinBySignedInContinuedThreeDay                          = 50
	RewardCoinBySignedInContinuedFourDay                           = 80
	RewardCoinBySignedInContinuedFiveDay                           = 120
	RewardCoinBySignedInContinuedSixDay                            = 200
	RewardCoinBySignedInContinuedSevenDay                          = 300
)

var RewardCoinBySignedInNextMap = map[RewardCoinBySignedInType]RewardCoinBySignedInType{
	RewardCoinBySignedInContinuedOneDay:   RewardCoinBySignedInContinuedTwoDay,
	RewardCoinBySignedInContinuedTwoDay:   RewardCoinBySignedInContinuedThreeDay,
	RewardCoinBySignedInContinuedThreeDay: RewardCoinBySignedInContinuedFourDay,
	RewardCoinBySignedInContinuedFourDay:  RewardCoinBySignedInContinuedFiveDay,
	RewardCoinBySignedInContinuedFiveDay:  RewardCoinBySignedInContinuedSixDay,
	RewardCoinBySignedInContinuedSixDay:   RewardCoinBySignedInContinuedSevenDay,
	RewardCoinBySignedInContinuedSevenDay: RewardCoinBySignedInContinuedOneDay,
}

var RewardCoinBySignInContinueDayMap = map[RewardCoinBySignedInType]int{
	RewardCoinBySignedInContinuedOneDay:   1,
	RewardCoinBySignedInContinuedTwoDay:   2,
	RewardCoinBySignedInContinuedThreeDay: 3,
	RewardCoinBySignedInContinuedFourDay:  4,
	RewardCoinBySignedInContinuedFiveDay:  5,
	RewardCoinBySignedInContinuedSixDay:   6,
	RewardCoinBySignedInContinuedSevenDay: 7,
}
