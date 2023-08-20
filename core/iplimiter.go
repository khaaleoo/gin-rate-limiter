package core

import (
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	ipLimiterIns map[string]*IPLimiter = make(map[string]*IPLimiter)
)

type IPLimiter struct {
	RateLimiter
}

func newIPLimiter(key string, option RateLimiterOption) *IPLimiter {
	mu.Lock()
	defer mu.Unlock()

	if ipLimiterIns[key] != nil {
		return ipLimiterIns[key]
	}

	instance := &IPLimiter{
		RateLimiter: RateLimiter{
			RateLimiterType: IPRateLimiter,
			Key:             key,
			Option:          option,
			Items:           make(map[string]*RateLimiterItem),
		},
	}
	ipLimiterIns[key] = instance
	return instance
}

type IPLimiterItem struct {
	RateLimiterItem
}

func (m *IPLimiter) newItem(ip string) *RateLimiterItem {
	if _, exist := ipLimiterIns[m.Key].Items[ip]; exist {
		ipLimiterIns[m.Key].Items[ip] = nil
	}

	instance := RateLimiterItem{
		Key:        ip,
		Limiter:    rate.NewLimiter(m.Option.Limit, m.Option.Burst),
		LastSeenAt: time.Now(),
	}
	ipLimiterIns[m.Key].setItem(ip, &instance)

	return &instance
}

func (m *IPLimiter) GetItem(ctx *gin.Context) (*RateLimiterItem, error) {
	mu.Lock()
	defer mu.Unlock()

	ip := ctx.ClientIP()
	item, exist := ipLimiterIns[m.Key].Items[ip]
	if !exist {
		return m.newItem(ip), nil
	}

	if time.Since(item.LastSeenAt) > ipLimiterIns[m.Key].Option.Len {
		return m.newItem(ip), nil
	}

	ipLimiterIns[m.Key].Items[ip].LastSeenAt = time.Now()
	return item, nil
}

func (m *IPLimiter) setItem(ip string, item *RateLimiterItem) error {
	if m.Items == nil {
		m.Items = make(map[string]*RateLimiterItem)
	}

	m.Items[ip] = item

	return nil
}
