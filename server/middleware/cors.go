package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type CorsHandler func() fiber.Handler

func NewCors() CorsHandler {
	return func() fiber.Handler {
		return cors.New(cors.Config{
			AllowCredentials: true,
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
			AllowOriginsFunc: func(origin string) bool {
				return true
			},
		})
	}
}
