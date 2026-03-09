package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = map[string]*visitor{}
	mu       sync.Mutex
)

func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	go cleanupVisitors()
	return func(c *gin.Context) {
		limiter := getVisitor(c.ClientIP(), requestsPerMinute)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

func getVisitor(ip string, requestsPerMinute int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	if v, ok := visitors[ip]; ok {
		v.lastSeen = time.Now()
		return v.limiter
	}
	limiter := rate.NewLimiter(rate.Every(time.Minute/time.Duration(requestsPerMinute)), requestsPerMinute)
	visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
	return limiter
}

func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}
