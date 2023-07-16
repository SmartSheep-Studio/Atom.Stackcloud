package controllers

import (
	"code.smartsheep.studio/atom/neutron/http/context"
	"code.smartsheep.studio/atom/stackcloud/datasource/models"
	"code.smartsheep.studio/atom/stackcloud/http/middleware"
	"code.smartsheep.studio/atom/stackcloud/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FunctionController struct {
	db       *gorm.DB
	executor *services.FunctionService
	auth     middleware.AuthHandler
}

func NewFunctionController(db *gorm.DB, executor *services.FunctionService, auth middleware.AuthHandler) *FunctionController {
	return &FunctionController{db, executor, auth}
}

func (ctrl *FunctionController) Map(router *context.App) {
	router.Get("/api/apps/:app/functions", ctrl.auth(true, "records.function.read", "stackcloud.function.read"), ctrl.list)
	router.Get("/api/apps/:app/functions/:function", ctrl.auth(true, "records.function.read", "stackcloud.function.read"), ctrl.get)
	router.Post("/api/apps/:app/functions", ctrl.auth(true, "records.function.create", "stackcloud.function.create"), ctrl.create)
	router.Put("/api/apps/:app/functions/:function", ctrl.auth(true, "records.function.update", "stackcloud.function.update"), ctrl.update)
	router.Delete("/api/apps/:app/functions/:function", ctrl.auth(true, "records.function.delete", "stackcloud.function.delete"), ctrl.delete)
	router.Post("/api/apps/:app/functions/:function/call", ctrl.auth(true, "records.function.call", "stackcloud.function.call"), ctrl.call)
}

func (ctrl *FunctionController) list(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var functions []models.CloudFunction
	if err := ctrl.db.Where("app_id = ?", app.ID).Order("created_at desc").Find(&functions).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(functions)
	}
}

func (ctrl *FunctionController) get(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var function models.CloudFunction
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("function"), app.ID).First(&function).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(function)
	}
}

func (ctrl *FunctionController) create(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Script      string   `json:"script" validate:"required"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	function := models.CloudFunction{
		Slug:        req.Slug,
		Name:        req.Name,
		Script:      req.Script,
		Description: req.Description,
		Tags:        datatypes.NewJSONSlice(req.Tags),
		AppID:       app.ID,
	}

	if err := ctrl.db.Save(&function).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(function)
	}
}

func (ctrl *FunctionController) update(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Script      string   `json:"script" validate:"required"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	var function models.CloudFunction
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("function"), app.ID).First(&function).Error; err != nil {
		return c.DbError(err)
	}

	function.Slug = req.Slug
	function.Name = req.Name
	function.Script = req.Script
	function.Description = req.Description
	function.Tags = datatypes.NewJSONSlice(req.Tags)

	if err := ctrl.db.Save(&function).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(function)
	}
}

func (ctrl *FunctionController) delete(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var function models.CloudFunction
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("function"), app.ID).First(&function).Error; err != nil {
		return c.DbError(err)
	}

	if err := ctrl.db.Delete(&function).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func (ctrl *FunctionController) call(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}

	var app models.App
	if err := ctrl.db.Where("slug = ?", c.Params("app")).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var function models.CloudFunction
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("function"), app.ID).First(&function).Error; err != nil {
		return c.DbError(err)
	}

	return ctrl.executor.HandleRequest(function, c)
}
