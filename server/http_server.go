package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"benkyou/model"
	"benkyou/service"
	"benkyou/utils"

	"github.com/labstack/echo"
)

func RunServer(ctx context.Context, kanjiService *service.Service, port string) {
	cfgBucketSize := os.Getenv("RATE_LIMITER_BUCKET_SIZE")
	rateLimiterBucketSize, err := strconv.ParseInt(cfgBucketSize, 10, 32)
	if err != nil {
		log.Fatal("fail parse rate limiter bucket size")
	}

	cfgTickInterval := os.Getenv("RATE_LIMITER_TICK_INTERVAL")
	tTickInterval, err := strconv.ParseInt(cfgTickInterval, 10, 32)
	if err != nil {
		log.Fatal("fail parse rate limiter tick interval")
	}

	rateLimiterTickInterval := time.Millisecond * time.Duration(tTickInterval)

	e := echo.New()
	e.Use(rateLimiterMiddleware(ctx, int(rateLimiterBucketSize), rateLimiterTickInterval))
	e.GET("/get-question", getQuestion(kanjiService))

	if err := e.Start(port); err != nil {
		log.Fatal(err)
	}
}

func rateLimiterMiddleware(ctx context.Context, bucketSize int, tickInterval time.Duration) echo.MiddlewareFunc {
	rl := utils.NewLimiter(ctx, bucketSize, tickInterval)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !rl.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"message": "too many requests",
				})
			}

			return next(c)
		}
	}
}

func getQuestion(kanjiSvc *service.Service) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		level := ctx.QueryParam("level")
		if level == "" {
			level = "n5"
		}

		wordType := ctx.QueryParam("wordType")
		if wordType == "" {
			wordType = model.WordTypeKeiyoushi
		}

		question, err := kanjiSvc.GetQuestionByLevel(level, wordType)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, question)
	}
}
