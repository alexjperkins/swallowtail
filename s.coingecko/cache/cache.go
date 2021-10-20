package cache

// CoingeckoCache ...
type CoingeckoCache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Close()
}
