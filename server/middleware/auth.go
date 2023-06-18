package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"repo.smartsheep.studio/atom/nucleus/datasource/models"
	"repo.smartsheep.studio/atom/nucleus/toolbox"
)

var conn *toolbox.ExternalServiceConnection

type AuthHandler func(force bool, perms ...string) fiber.Handler

type AuthConfig struct {
	Next        func(c *fiber.Ctx) bool
	LookupToken string
}

func NewAuth(cycle fx.Lifecycle, c *toolbox.ExternalServiceConnection) AuthHandler {
	conn = c

	cfg := AuthConfig{
		Next:        nil,
		LookupToken: "header: Authorization, query: token, cookie: lineup_authorization",
	}

	return func(force bool, perms ...string) fiber.Handler {
		return func(c *fiber.Ctx) error {
			if cfg.Next != nil && cfg.Next(c) {
				return c.Next()
			}

			u, err := LookupAuthToken(c, strings.Split(cfg.LookupToken, ","))
			if err != nil && force {
				return fiber.NewError(fiber.StatusUnauthorized, err.Error())
			} else {
				if err == nil {
					if err := conn.HasUserPermissions(u, perms...); err != nil {
						return fiber.NewError(fiber.StatusForbidden, err.Error())
					}
				}

				c.Locals("principal-ok", err == nil)

				c.Locals("principal", u)
			}

			return c.Next()
		}
	}
}

func LookupAuthToken(c *fiber.Ctx, args []string) (models.User, error) {
	var str string
	for _, arg := range args {
		parts := strings.Split(strings.TrimSpace(arg), ":")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])

		switch k {
		case "header":
			if len(c.GetReqHeaders()[v]) > 0 {
				str = strings.TrimSpace(strings.ReplaceAll(c.GetReqHeaders()[v], "Bearer", ""))
			}
		case "query":
			if len(c.Query(v)) > 0 {
				str = c.Query(v)
			}
		case "cookie":
			if len(c.Cookies(v)) > 0 {
				str = c.Cookies(v)
			}
		}
	}

	if len(str) == 0 {
		return models.User{}, fmt.Errorf("could not found any token string from context")
	}

	return conn.GetPrincipal(str, true)
}
