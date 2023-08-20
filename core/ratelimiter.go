package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var mu sync.Mutex

type IRateLimiter interface {
	newItem(itemKey string) *RateLimiterItem
	setItem(key string, item *RateLimiterItem) error
	GetItem(ctx *gin.Context) (*RateLimiterItem, error)
}

type RateLimiter struct {
	RateLimiterType RateLimiterType
	Key             string
	Option          RateLimiterOption
	Items           map[string]*RateLimiterItem
}

type RateLimiterOption struct {
	Limit rate.Limit
	Burst int
	Len   time.Duration
}

type RateLimiterItem struct {
	Key        string
	Limiter    *rate.Limiter
	LastSeenAt time.Time
}

func RequireRateLimiter(rateLimiters ...RateLimiter) func(*gin.Context) {
	return func(c *gin.Context) {
		for _, rateLimiter := range rateLimiters {
			instance, err := GetRateLimiterInstance(rateLimiter.RateLimiterType, rateLimiter.Key, rateLimiter.Option)
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{
					"message": err.Error(),
				})
				return
			}

			item, err := instance.GetItem(c)
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{
					"message": err.Error(),
				})
				return
			}

			if !item.Limiter.Allow() {
				c.AbortWithStatusJSON(429, gin.H{
					"message": "Too many requests",
				})
				return
			}
		}
		c.Next()
	}
}

func GetRateLimiterInstance(rateLimiterType RateLimiterType, key string, option RateLimiterOption) (IRateLimiter, error) {
	switch rateLimiterType {
	case IPRateLimiter:
		return newIPLimiter(key, option), nil
	case JWTRateLimiter:
		return nil, fmt.Errorf("JWTRateLimiter is not implemented yet")
	default:
		return nil, fmt.Errorf("rateLimiterType %v is not supported", rateLimiterType)
	}
}
