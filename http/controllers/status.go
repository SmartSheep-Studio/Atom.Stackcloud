package controllers

import (
	ctx "code.smartsheep.studio/atom/neutron/http/context"
	"code.smartsheep.studio/atom/neutron/toolbox"
)

type StatusController struct {
	conn *toolbox.ExternalServiceConnection
}

func NewStatusController(conn *toolbox.ExternalServiceConnection) *StatusController {
	return &StatusController{conn}
}

func (ctrl *StatusController) Map(router *ctx.App) {
	router.Get("/api/info", ctrl.configure)
}

func (ctrl *StatusController) configure(c *ctx.Ctx) error {
	return c.P.JSON(ctrl.conn)
}
