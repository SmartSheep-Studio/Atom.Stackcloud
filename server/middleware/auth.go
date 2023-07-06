package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"repo.smartsheep.studio/atom/matrix/datasource/models"
	tmodels "repo.smartsheep.studio/atom/nucleus/datasource/models"
	"repo.smartsheep.studio/atom/nucleus/toolbox"
)

var conn *toolbox.ExternalServiceConnection

type AuthHandler func(force bool, perms ...string) fiber.Handler

type AuthConfig struct {
	Next        func(c *fiber.Ctx) bool
	LookupToken string
}

func NewAuth(cycle fx.Lifecycle, db *gorm.DB, c *toolbox.ExternalServiceConnection) AuthHandler {
	conn = c

	cfg := AuthConfig{
		Next:        nil,
		LookupToken: "header: Authorization, query: token, cookie: authorization",
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

					var account *models.MatrixAccount
					if err := db.Where("user_id = ?", u.ID).First(&account).Error; err != nil {
						if errors.Is(gorm.ErrRecordNotFound, err) {
							account = &models.MatrixAccount{
								Nickname: u.Nickname,
								UserID:   u.ID,
							}

							if err := db.Save(&account).Error; err != nil {
								return fiber.NewError(fiber.StatusInternalServerError, err.Error())
							}
						} else {
							return fiber.NewError(fiber.StatusInternalServerError, err.Error())
						}
					}

					c.Locals("matrix-id", account)
				}

				c.Locals("principal-ok", err == nil)

				c.Locals("principal", u)
			}

			return c.Next()
		}
	}
}

func LookupAuthToken(c *fiber.Ctx, args []string) (tmodels.User, error) {
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
		return tmodels.User{}, fmt.Errorf("could not found any token string from context")
	}

	resp, err := conn.GetPrincipal(str)
	if err != nil {
		return tmodels.User{}, err
	} else {
		return resp.User, nil
	}
}
