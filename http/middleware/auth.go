package middleware

import (
	ctx "code.smartsheep.studio/atom/neutron/http/context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"

	"code.smartsheep.studio/atom/matrix/datasource/models"
	tmodels "code.smartsheep.studio/atom/neutron/datasource/models"
	"code.smartsheep.studio/atom/neutron/toolbox"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var conn *toolbox.ExternalServiceConnection

type AuthHandler func(force bool, perms ...string) ctx.Handler

type AuthConfig struct {
	Next        func(c *ctx.Ctx) bool
	LookupToken string
}

func NewAuth(cycle fx.Lifecycle, db *gorm.DB, c *toolbox.ExternalServiceConnection) AuthHandler {
	conn = c

	cfg := AuthConfig{
		Next:        nil,
		LookupToken: "header: Authorization, query: token, cookie: authorization",
	}

	return func(force bool, perms ...string) ctx.Handler {
		return func(c *ctx.Ctx) error {
			if cfg.Next != nil && cfg.Next(c) {
				return c.P.Next()
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

					c.P.Locals("matrix-id", account)
				}

				c.P.Locals("principal-ok", err == nil)

				c.P.Locals("principal", u)
			}

			return c.P.Next()
		}
	}
}

func LookupAuthToken(c *ctx.Ctx, args []string) (tmodels.User, error) {
	var str string
	for _, arg := range args {
		parts := strings.Split(strings.TrimSpace(arg), ":")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])

		switch k {
		case "header":
			if len(c.P.GetReqHeaders()[v]) > 0 {
				str = strings.TrimSpace(strings.ReplaceAll(c.P.GetReqHeaders()[v], "Bearer", ""))
			}
		case "query":
			if len(c.P.Query(v)) > 0 {
				str = c.P.Query(v)
			}
		case "cookie":
			if len(c.P.Cookies(v)) > 0 {
				str = c.P.Cookies(v)
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
