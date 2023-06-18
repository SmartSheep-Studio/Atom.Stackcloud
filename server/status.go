package server

import (
	"github.com/gofiber/fiber/v2"
	"repo.smartsheep.studio/atom/nucleus/toolbox"
)

type StatusController struct {
	conn *toolbox.ExternalServiceConnection
}

func NewStatusController(conn *toolbox.ExternalServiceConnection) *StatusController {
	return &StatusController{conn}
}

func (ctrl *StatusController) Map(router *fiber.App) {
	router.Get("/api/info", ctrl.configure)
}

func (ctrl *StatusController) configure(c *fiber.Ctx) error {
	return c.JSON(ctrl.conn)
}
