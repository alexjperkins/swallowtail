package cache

// CoingeckoCache ...
type CoingeckoCache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}
