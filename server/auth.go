package server

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"repo.smartsheep.studio/atom/nucleus/datasource/models"
	"repo.smartsheep.studio/atom/nucleus/toolbox"
)

type AuthController struct {
	conn *toolbox.ExternalServiceConnection
}

func NewAuthController(conn *toolbox.ExternalServiceConnection) *AuthController {
	return &AuthController{conn}
}

func (ctrl *AuthController) Map(router *fiber.App) {
	router.Get("/api/auth/request", ctrl.request)
	router.Get("/api/auth/callback", ctrl.callback)
}

func (ctrl *AuthController) getOauth() models.OauthClient {
	raw, _ := json.Marshal(ctrl.conn.Additional["oauth"])

	var client models.OauthClient
	_ = json.Unmarshal(raw, &client)

	return client
}

func (ctrl *AuthController) request(c *fiber.Ctx) error {
	oauth := ctrl.getOauth()

	return c.Redirect(ctrl.conn.GetConnectURL(
		strconv.Itoa(int(oauth.ID)),
		fmt.Sprintf("%s/api/auth/callback", viper.GetString("general.base_url")),
	), fiber.StatusFound)
}

func (ctrl *AuthController) callback(c *fiber.Ctx) error {
	oauth := ctrl.getOauth()

	code := c.Query("code")

	token, _, err := ctrl.conn.ExchangeAccessToken(
		code,
		strconv.Itoa(int(oauth.ID)),
		oauth.Secret,
		fmt.Sprintf("%s/api/auth/callback", viper.GetString("general.base_url")),
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		c.Cookie(&fiber.Cookie{
			Path:   "/",
			Name:   "lineup_authorization",
			Value:  token,
			MaxAge: int(viper.GetDuration("security.sessions_alive_duration").Seconds()),
		})

		return c.Redirect(viper.GetString("general.base_url"), fiber.StatusFound)
	}
}
