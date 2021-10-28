package client

import "net/http"

type bitfinexRateLimiter struct{}

func (b *bitfinexRateLimiter) RefreshWait(header http.Header, statusCode int) {}
func (b *bitfinexRateLimiter) Wait()                                          {}
