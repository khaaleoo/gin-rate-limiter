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

// This test case is to test the rate limiter for maximum number of requests per second
// The rate limiter is set to maximum 5 requests per second for each IP
// The test case will send 20 requests to the server and check if the server returns 429 error for 15 requests
func TestMaximumRequestsInAPeriod(t *testing.T) {
	r := SetUpRouter()
	WindowCapacity := 5
	WindowLen := 1 * time.Second
	rateLimiterMiddleware := ratelimiter.RequireRateLimiter(ratelimiter.RateLimiter{
		RateLimiterType: ratelimiter.IPRateLimiter,
		Key:             "iplimiter_maximum_requests_for_ip_test",
		Option: ratelimiter.RateLimiterOption{
			Limit: 1,
			Burst: WindowCapacity,
			Len:   WindowLen,
		},
	})

	r.GET("/me", rateLimiterMiddleware, CommonHandler)

	// Fake Requests
	successCount := 0
	errorCount := 0
	numRequests := 20
	requests := make([]int, numRequests)

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

	wg.Wait()

	timeEnd := time.Now()

	t.Log("success count: ", successCount)
	t.Log("error count: ", errorCount)
	t.Log("time elapsed: ", timeEnd.Sub(timeStart))
	assert.Equal(t, numRequests-WindowCapacity, errorCount, fmt.Sprintf("error count should be %d", numRequests-WindowCapacity))
}
