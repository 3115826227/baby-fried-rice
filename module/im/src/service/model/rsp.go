package model

type RspSuccess struct {
	Code int `json:"code"`
}

type RspOkResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RespSuccessData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RspFriendCategory struct {
	Name    string      `json:"name"`
	Friends []RspFriend `json:"friends"`
}

type RspFriend struct {
	Id       string `json:"id"`
	Friend   string `json:"friend"`
	Username string `json:"username"`
	Remark   string `json:"remark"`
	HeadImg  string `json:"head_img"`
}
