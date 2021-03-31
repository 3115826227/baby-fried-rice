package constant

type SendPrivateMessageType int

const (
	SendPerson SendPrivateMessageType = 1
	SendGroup                         = 2
	SendGlobal                        = 3
)
