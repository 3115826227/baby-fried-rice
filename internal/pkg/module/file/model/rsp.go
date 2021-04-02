package model

type FileUploadResp struct {
	ID         string `json:"id"`
	Path       string `json:"path"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Origin     string `json:"origin"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type FileResp struct {
	ID             string `json:"id"`
	Path           string `json:"path"`
	Name           string `json:"name"`
	Size           int64  `json:"size"`
	Origin         string `json:"origin"`
	Type           int    `json:"type"`
	PermissionType int    `json:"permission_type"`
	CreateTime     string `json:"create_time"`
	UpdateTime     string `json:"update_time"`
}
