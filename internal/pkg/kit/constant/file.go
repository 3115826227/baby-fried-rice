package constant

type FilePermissionType int

const (
	PublicFile  FilePermissionType = 1
	ProtectFile                    = 2
	PrivateFile                    = 3
)

type FileMode string

const (
	LocalFileMode       FileMode = "local"
	QiNiuYunOssFileMode          = "qiniuyun"
)
