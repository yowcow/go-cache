package cache

type Cacher interface {
	MaxSize() int64
	CurrentSize() int64
	AllKeys() []string
	AllKeysReversed() []string
	Set(string, interface{}) error
	Get(string) (interface{}, error)
	Delete(string) error
}
