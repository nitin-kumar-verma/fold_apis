package middlewares

import (
	"fold/internal/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var supportedLanguages map[string]string = map[string]string{
	"en-IN": "en-IN",
}

func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//Change below to restrict network calls from specific urls
		c.Set("Access-Control-Allow-Origin", "*")
		//remove methods which are not required
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		//block on the basis of allowed headers
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		return c.Next()
	}
}

func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		//JWT validation logic, connect to keycloak and authenticate
		return c.Next()
	}
}

// This method initialized a new logger instance and sets req. Id
func SetRequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			var err error
			requestID, err = logger.GenerateRequestID()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}
			c.Set("X-Request-ID", requestID)
		}

		// Set the request ID in the context
		c.Locals("RequestID", requestID)

		// Create a new logger instance for the request
		logger := logger.InitLogger(requestID)

		// Set the logger in the context
		c.Locals("Logger", logger)

		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", zap.Any("error", err))
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
			}
		}()

		logger.Info("Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		)

		return c.Next()
	}
}
