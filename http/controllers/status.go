package controllers

import (
	"code.smartsheep.studio/atom/neutron/http/context"
	"code.smartsheep.studio/atom/neutron/toolbox"
	"github.com/gofiber/fiber/v2"
)

type StatusController struct {
	conn *toolbox.ExternalServiceConnection
}

func NewStatusController(conn *toolbox.ExternalServiceConnection) *StatusController {
	return &StatusController{conn}
}

func (ctrl *StatusController) Map(router *context.App) {
	router.Get("/api/info", ctrl.configure)
}

func (ctrl *StatusController) configure(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	return c.JSON(ctrl.conn)
}
