package requests

type ReqAddSession struct {
	SessionType        int32    `json:"session_type"`
	JoinPermissionType int32    `json:"join_permission_type"`
	Name               string   `json:"name"`
	Joins              []string `json:"joins"`
}
