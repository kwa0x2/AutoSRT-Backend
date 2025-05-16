package middleware

import (
	"net/http"

	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/gin-gonic/gin"
)

func PaddleWebhookVerifier(secretKey string) gin.HandlerFunc {
	verifier := paddle.NewWebhookVerifier(secretKey)

	return func(c *gin.Context) {
		handler := verifier.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		}))

		handler.ServeHTTP(c.Writer, c.Request)

		if c.IsAborted() {
			return
		}
	}
}
