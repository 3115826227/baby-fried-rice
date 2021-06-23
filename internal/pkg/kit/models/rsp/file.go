package rsp

type File struct {
	ID         string `json:"id"`
	Origin     string `json:"origin"`
	Name       string `json:"name"`
	DownUrl    string `json:"down_url"`
	Size       int64  `json:"size"`
	UploadTime int64  `json:"upload_time"`
	StorageDay int    `json:"storage_day"`
}

type FileUploadResp struct {
	File
}

type FileQueryResp struct {
	List []File `json:"list"`
}
