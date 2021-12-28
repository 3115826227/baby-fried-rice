package interfaces

type ConfigCenter interface {
	// 配置鉴权路由
	AddAuthPath(path, name string) error
	// 判断鉴权路由是否配置过
	IsAuthPathConfig(path string) (bool, error)
	// 获取所有鉴权路由
	GetAuthPath() (map[string]string, error)
	// 删除鉴权路由
	DelAuthPath(path string) error
}
