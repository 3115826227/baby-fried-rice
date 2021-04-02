package interfaces

type Cache interface {
	Add(key string, value string) error
	Get(key string) (string, error)
	Del(key string) error
}
