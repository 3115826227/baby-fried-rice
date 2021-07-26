package handle

const (
	SuccessCode = 0
)

type RspSuccess struct {
	Code int `json:"code"`
}

type RspOkResponse struct {
	Code    int    `json:"code"`
	Message string `json:"smsDao"`
}

type RespSuccessData struct {
	Code    int         `json:"code"`
	Message string      `json:"smsDao"`
	Data    interface{} `json:"data"`
}
