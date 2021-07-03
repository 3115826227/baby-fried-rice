package rsp

type User struct {
	AccountID  string `json:"account_id"`
	Username   string `json:"username"`
	HeadImgUrl string `json:"head_img_url"`
	Remark     string `json:"remark"`
}

type UserDataResp struct {
	UserId    string `json:"user_id"`
	Username  string `json:"username"`
	LoginName string `json:"login_name"`
}

type LoginResult struct {
	UserInfo UserDataResp `json:"user_info"`
	Token    string       `json:"token"`
}

type UserDetailResp struct {
	AccountId  string `json:"account_id"`
	Describe   string `json:"describe"`
	HeadImgUrl string `json:"head_img_url"`
	Username   string `json:"username"`
	SchoolId   string `json:"school_id"`
	Gender     bool   `json:"gender"`
	Age        int64  `json:"age"`
	Phone      string `json:"phone"`
}
