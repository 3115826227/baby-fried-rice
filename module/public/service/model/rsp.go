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

type RspGrade struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RspSubject struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RspStreet struct {
	Street string `json:"street"`
}

type RspLocal struct {
	Local   string      `json:"local"`
	Code    string      `json:"code"`
	Streets []RspStreet `json:"streets"`
}

type RspCity struct {
	City   string     `json:"city"`
	Code   string     `json:"code"`
	Locals []RspLocal `json:"locals"`
}

type RspArea struct {
	Province string    `json:"province"`
	Code     string    `json:"code"`
	Cities   []RspCity `json:"cities"`
}
