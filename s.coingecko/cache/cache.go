package cache

// CoingeckoCache ...
type CoingeckoCache interface {
	Get(key string) (interface{}, bool, error)
	Set(key string, value interface{}) error
}
