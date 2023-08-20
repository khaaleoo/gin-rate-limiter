package example

import (
	"time"

	"github.com/gin-gonic/gin"
	ratelimiter "github.com/khaaleoo/gin-rate-limiter/core"
)

func Example() {
	r := gin.Default()

	// Create an IP rate limiter middleware
	rateLimiterMiddleware := ratelimiter.RequireRateLimiter(ratelimiter.RateLimiter{
		RateLimiterType: ratelimiter.IPRateLimiter,
		Key:             "iplimiter_maximum_requests_for_ip_test",
		Option: ratelimiter.RateLimiterOption{
			Limit: 1,
			Burst: 1,
			Len:   1 * time.Second,
		},
	})

	// Apply rate limiter middleware to a route
	r.GET("/limited-route", rateLimiterMiddleware, func(c *gin.Context) {
		c.String(200, "Hello, rate-limited world!")
	})

	r.Run(":8080")
}
