package requests

type PasswordLoginReq struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`   // 密码
	Ip        string `json:"ip"`
}

type UserRegisterReq struct {
	PasswordLoginReq
	Username string `json:"username" binding:"required"` //昵称
	Gender   bool   `json:"gender" binding:"required"`   //性别
	Phone    string `json:"phone" binding:"required"`    //手机号
}

// 管理平台添加用户
type AddUserReq struct {
	UserRegisterReq
	// 是否为官方用户
	IsOfficial bool `json:"is_official"`
}

type UserDetailUpdateReq struct {
	HeadImgUrl string `json:"head_img_url"`
	Describe   string `json:"describe"`
	Username   string `json:"username"`
	Gender     bool   `json:"gender"`
	Phone      string `json:"phone"`
	Age        int64  `json:"age"`
}

type UserPwdUpdateReq struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

type UserCoinGiveawayReq struct {
	Coin int64    `json:"coin"`
	Ids  []string `json:"ids"`
}
