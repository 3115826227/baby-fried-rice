package interfaces

type Cache interface {
	Add(key string, value interface{}) error
	Get(key string) (string, error)
	Del(key string) error
}
