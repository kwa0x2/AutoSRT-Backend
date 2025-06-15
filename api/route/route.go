package route

import (
	"net/http"

	"github.com/PaddleHQ/paddle-go-sdk/v3"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/middleware"
	"github.com/kwa0x2/AutoSRT-Backend/config"
	"github.com/resend/resend-go/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Setup(env *config.Env, db *mongo.Database, dynamodb *dynamodb.Client, router *gin.Engine, resendClient *resend.Client, s3Client *s3.Client, lambdaClient *lambda.Client, paddleSDK *paddle.SDK) {
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404 made by kwa -> https://github.com/kwa0x2")
	})

	limiter := middleware.NewRateLimiter()

	groupRouter := router.Group("/api/v1")
	groupRouter.Use(limiter.RateLimitMiddleware())

	NewAuthRoute(env, groupRouter, db, dynamodb, resendClient, paddleSDK)
	NewUserRoute(groupRouter, db, dynamodb)
	NewSRTRoute(groupRouter, s3Client, lambdaClient, env.AWSS3BucketName, env.AWSLambdaFuncName, db, dynamodb)
	NewUsageRoute(groupRouter, db, dynamodb)
	NewContactRoute(env, groupRouter, db, resendClient)
	NewPaddleRoutes(env, groupRouter, paddleSDK, db, dynamodb)
}
