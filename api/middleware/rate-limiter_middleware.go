package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limits   map[string]EndpointLimit
}

type EndpointLimit struct {
	limit  int
	window time.Duration
}

var defaultEndpointLimits = map[string]EndpointLimit{
	// Auth endpoints
	"GET/api/v1/auth/google/login":             {limit: 10, window: time.Minute},
	"GET/api/v1/auth/github/login":             {limit: 10, window: time.Minute},
	"POST/api/v1/auth/credentials/login":       {limit: 10, window: time.Minute},
	"POST/api/v1/auth/register":                {limit: 10, window: time.Minute},
	"GET/api/v1/auth/logout":                   {limit: 15, window: time.Minute},
	"POST/api/v1/auth/otp/send":                {limit: 2, window: time.Minute},
	"POST/api/v1/auth/account/password/forgot": {limit: 5, window: time.Minute},
	"PUT/api/v1/auth/account/password/reset":   {limit: 5, window: time.Minute},
	"GET/api/v1/auth/account/delete/request":   {limit: 5, window: time.Minute},
	"DELETE/api/v1/auth/account":               {limit: 5, window: time.Minute},

	// User endpoints
	"GET/api/v1/user/me":                   {limit: 500, window: time.Minute},
	"HEAD/api/v1/user/exists/email/:email": {limit: 20, window: time.Minute},
	"HEAD/api/v1/user/exists/phone/:phone": {limit: 20, window: time.Minute},
	// SRT endpoints
	"POST/api/v1/srt":          {limit: 10, window: time.Minute},
	"GET/api/v1/srt/histories": {limit: 100, window: time.Minute},
	// Usage endpoint
	"GET/api/v1/usage": {limit: 500, window: time.Minute},

	// Contact endpoint
	"POST/api/v1/contact": {limit: 5, window: time.Minute},

	// Paddle endpoints
	"POST/api/v1/paddle/webhook":        {limit: 200, window: time.Minute},
	"GET/api/v1/paddle/customer-portal": {limit: 50, window: time.Minute},
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limits:   defaultEndpointLimits,
	}
}

func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		endpointKey := ctx.Request.Method + ctx.FullPath()
		clientKey := ctx.ClientIP() + ":" + endpointKey

		limitConfig, exists := rl.limits[endpointKey]
		if !exists {
			ctx.Next()
			return
		}

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		var validRequests []time.Time

		for _, t := range rl.requests[clientKey] {
			if now.Sub(t) < limitConfig.window {
				validRequests = append(validRequests, t)
			}
		}
		rl.requests[clientKey] = validRequests

		if len(rl.requests[clientKey]) >= limitConfig.limit {
			ctx.JSON(http.StatusTooManyRequests, utils.NewMessageResponse("Rate limit exceeded. Please try again later or contact support."))
			ctx.Abort()
			return
		}

		rl.requests[clientKey] = append(rl.requests[clientKey], now)

		ctx.Next()
	}
}
