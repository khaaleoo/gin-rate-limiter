package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	ratelimiter "github.com/khaaleoo/gin-rate-limiter/core"
	"github.com/stretchr/testify/assert"
)

// This test case is to test the rate limiter for different routers using the same rate limiter middleware
// Same rate limiter middleware means the same type of rate limiter (IPRateLimiter) and the same key ("ip")
// The rate limiter is set to maximum 5 requests per second for each IP
// We will make 5 requests to /ping and 5 requests to /me, summing up to 10 requests
// The test case will pass if the server returns 429 error for 5 requests
func TestMaximumRequestInDifferentRoutesUsingSameMiddleware(t *testing.T) {
	r := SetUpRouter()
	WindowCapacity := 5
	WindowLen := 1 * time.Second
	rateLimiterMiddleware := ratelimiter.RequireRateLimiter(ratelimiter.RateLimiter{
		RateLimiterType: ratelimiter.IPRateLimiter,
		Key:             "iplimiter_different_route_same_middleware_test",
		Option: ratelimiter.RateLimiterOption{
			Limit: 1,
			Burst: WindowCapacity,
			Len:   WindowLen,
		},
	})

	r.GET("/ping", rateLimiterMiddleware, CommonHandler)
	r.GET("/me", rateLimiterMiddleware, CommonHandler)

	// Fake Requests
	successCount := 0
	errorCount := 0
	numRequestsPerRoute := 5
	numRoutes := 2
	requests := make([]int, numRequestsPerRoute)

	var wg sync.WaitGroup

	timeStart := time.Now()

	for i := range requests {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req, _ := http.NewRequest("GET", "/me", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code == http.StatusOK {
				successCount++
			} else if w.Code == http.StatusTooManyRequests {
				errorCount++
			}
		}(i)
	}

	for i := range requests {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req, _ := http.NewRequest("GET", "/ping", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code == http.StatusOK {
				successCount++
			} else if w.Code == http.StatusTooManyRequests {
				errorCount++
			}
		}(i)
	}

	wg.Wait()

	timeEnd := time.Now()

	t.Log("success count: ", successCount)
	t.Log("error count: ", errorCount)
	t.Log("time elapsed: ", timeEnd.Sub(timeStart))
	assert.Equal(t, numRequestsPerRoute*numRoutes-WindowCapacity, errorCount, fmt.Sprintf("error count should be %d", numRequestsPerRoute*numRoutes-WindowCapacity))
}
