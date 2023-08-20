package core

import (
	"time"

	"golang.org/x/time/rate"
)

const (
	IPLimiterWindowLen            = 1 * time.Second
	IPLimiterBurst                = 15
	IPLimiterLimit     rate.Limit = 1
)

const (
	IPRateLimiter RateLimiterType = iota
	JWTRateLimiter
)

type RateLimiterType int
