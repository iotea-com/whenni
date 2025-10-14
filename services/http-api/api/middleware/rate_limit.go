package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/config"
)

var redisStore = redis.New(redis.Config{
	Host:     config.VaultConf.RedisHost,
	Port:     config.VaultConf.RedisPort,
	Username: config.VaultConf.RedisUsername,
	Password: config.VaultConf.RedisPassword,
})

var userRateLimiterConfig = limiter.Config{
	KeyGenerator: func(ctx *fiber.Ctx) string {
		// 1. Use org ID from API key
		if apiKeyOrgId := ctx.Locals("api_key_org_id"); apiKeyOrgId != nil {
			return fmt.Sprintf("rate:org:api_key:%s", apiKeyOrgId)
		}

		// 2. Use user ID from JWT
		if userID := ctx.Locals("user_id"); userID != nil {
			return fmt.Sprintf("rate:user:%s", userID)
		}

		// 3. Fallback to IP
		return fmt.Sprintf("rate:ip_fallback:%s", ctx.IP())
	},
	Max:               200,
	Expiration:        1 * time.Minute,
	Storage:           redisStore,
	LimiterMiddleware: limiter.SlidingWindow{},
	LimitReached: func(c *fiber.Ctx) error {
		responseBody := ioteahttp.NewErrorResponse(
			[]any{"Too many requests"},
		)

		return c.Status(fiber.StatusTooManyRequests).JSON(responseBody)
	},
}

var ipRateLimiterConfig = limiter.Config{
	KeyGenerator: func(ctx *fiber.Ctx) string {
		return fmt.Sprintf("rate:ip:%s", ctx.IP())
	},
	Max:               100,
	Expiration:        1 * time.Minute,
	Storage:           redisStore,
	LimiterMiddleware: limiter.SlidingWindow{},
	LimitReached: func(ctx *fiber.Ctx) error {
		responseBody := ioteahttp.NewErrorResponse(
			[]any{"Too many requests"},
		)

		return ctx.Status(fiber.StatusTooManyRequests).JSON(responseBody)
	},
}

var UserRateLimiter = limiter.New(userRateLimiterConfig)
var IpRateLimiter = limiter.New(ipRateLimiterConfig)
