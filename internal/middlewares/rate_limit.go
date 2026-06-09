package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type clientInfo struct {
	Count     int64
	WindowEnd time.Time
}

var (
	mu          sync.Mutex
	clients     = make(map[string]*clientInfo)
	cleanupOnce sync.Once
)

func RateLimitMiddleware(maxRequests int64, window time.Duration) gin.HandlerFunc {

	cleanupOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(window * 2)
			defer ticker.Stop()
			for range ticker.C {
				now := time.Now()
				mu.Lock()
				for key, info := range clients {
					if now.After(info.WindowEnd) {
						delete(clients, key)
					}
				}
				mu.Unlock()
			}
		}()
	})

	return func(c *gin.Context) {

		// From authentication middleware
		userID := c.GetString("user_id")

		// Use IP as a fallback
		if userID == "" {
			userID = c.ClientIP()
		}

		now := time.Now()
		mu.Lock()
		info, exists := clients[userID]
		if !exists || now.After(info.WindowEnd) {
			clients[userID] = &clientInfo{
				Count:     1,
				WindowEnd: now.Add(window),
			}
			mu.Unlock()
			c.Next()
			return
		}

		if info.Count >= maxRequests {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "rate limit exceeded"})
			return
		}

		info.Count++
		mu.Unlock()
		c.Next()
	}
}
