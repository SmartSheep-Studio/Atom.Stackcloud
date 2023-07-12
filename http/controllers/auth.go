package controllers

import (
	ctx "code.smartsheep.studio/atom/neutron/http/context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"

	"code.smartsheep.studio/atom/neutron/datasource/models"
	"code.smartsheep.studio/atom/neutron/toolbox"
	"github.com/spf13/viper"
)

type AuthController struct {
	conn *toolbox.ExternalServiceConnection
}

func NewAuthController(conn *toolbox.ExternalServiceConnection) *AuthController {
	return &AuthController{conn}
}

func (ctrl *AuthController) Map(router *ctx.App) {
	router.Get("/api/auth/request", ctrl.request)
	router.Get("/api/auth/callback", ctrl.callback)
}

func (ctrl *AuthController) getOauth() models.OauthClient {
	raw, _ := json.Marshal(ctrl.conn.Additional["oauth"])

	var client models.OauthClient
	_ = json.Unmarshal(raw, &client)

	return client
}

func (ctrl *AuthController) request(c *ctx.Ctx) error {
	oauth := ctrl.getOauth()

	return c.P.Redirect(ctrl.conn.GetConnectURL(
		strconv.Itoa(int(oauth.ID)),
		fmt.Sprintf("%s/api/auth/callback", viper.GetString("base_url")),
	), fiber.StatusFound)
}

func (ctrl *AuthController) callback(c *ctx.Ctx) error {
	oauth := ctrl.getOauth()

	code := c.P.Query("code")

	token, _, err := ctrl.conn.ExchangeAccessToken(
		code,
		strconv.Itoa(int(oauth.ID)),
		oauth.Secret,
		fmt.Sprintf("%s/api/auth/callback", viper.GetString("base_url")),
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		c.P.Cookie(&fiber.Cookie{
			Path:     "/",
			Name:     "authorization",
			Value:    token,
			MaxAge:   int(viper.GetDuration("security.sessions_alive_duration").Seconds()),
			SameSite: "None",
			Secure:   true,
		})

		return c.P.Redirect(viper.GetString("base_url"), fiber.StatusFound)
	}
}
