package interfaces

type Cache interface {
	// kv操作
	Add(key string, value string) error
	Get(key string) (string, error)
	Del(key string) error
	// hash操作
	HSet(key, field string, value interface{}) error
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	HDel(key string, field ...string) error
	// 监控信息查询
	Info() (string, error)
}
