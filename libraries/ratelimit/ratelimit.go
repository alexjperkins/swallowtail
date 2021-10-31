package ratelimit

// RateLimiterOpts ...
type RateLimiterOpts struct {
	Metadata map[string]interface{}
}

// RateLimiter ...
type RateLimiter interface {
	Throttle()
	ThrottleWithOptions(*RateLimiterOpts)
}
