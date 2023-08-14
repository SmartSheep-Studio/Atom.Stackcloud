package controllers

import (
	"code.smartsheep.studio/atom/stackcloud/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/hypertext/middleware"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FunctionController struct {
	db         *gorm.DB
	executor   *services.FunctionService
	gatekeeper *middleware.AuthMiddleware
}

func NewFunctionController(db *gorm.DB, executor *services.FunctionService, gatekeeper *middleware.AuthMiddleware) *FunctionController {
	return &FunctionController{db, executor, gatekeeper}
}

func (ctrl *FunctionController) Map(router *fiber.App) {
	router.Get(
		"/api/apps/:app/functions",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:functions"), hyperutils.GenPerms("functions.read")),
		ctrl.list,
	)
	router.Get(
		"/api/apps/:app/functions/:function",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:functions"), hyperutils.GenPerms("functions.read")),
		ctrl.get,
	)
	router.Post(
		"/api/apps/:app/functions",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("create:functions"), hyperutils.GenPerms("functions.create")),
		ctrl.create,
	)
	router.Put(
		"/api/apps/:app/functions/:function",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("update:functions"), hyperutils.GenPerms("functions.update")),
		ctrl.update,
	)
	router.Delete(
		"/api/apps/:app/functions/:function",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("delete:functions"), hyperutils.GenPerms("functions.delete")),
		ctrl.delete,
	)
	router.Post(
		"/api/apps/:app/functions/:function/call",
		ctrl.gatekeeper.Fn(false, hyperutils.GenScope("call:functions"), hyperutils.GenPerms()),
		ctrl.call,
	)
}

func (ctrl *FunctionController) list(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var functions []models.CloudFunction
	if err := ctrl.db.Where("app_id = ?", app.ID).Order("created_at desc").Find(&functions).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(functions)
	}
}

func (ctrl *FunctionController) get(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var function models.CloudFunction
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("function"), app.ID).First(&function).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(function)
	}
}

func (ctrl *FunctionController) create(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Script      string   `json:"script" validate:"required"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
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
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(function)
	}
}

func (ctrl *FunctionController) update(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Script      string   `json:"script" validate:"required"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	var function models.CloudFunction
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("function"), app.ID).First(&function).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	function.Slug = req.Slug
	function.Name = req.Name
	function.Script = req.Script
	function.Description = req.Description
	function.Tags = datatypes.NewJSONSlice(req.Tags)

	if err := ctrl.db.Save(&function).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(function)
	}
}

func (ctrl *FunctionController) delete(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var function models.CloudFunction
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("function"), app.ID).First(&function).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	if err := ctrl.db.Delete(&function).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func (ctrl *FunctionController) call(c *fiber.Ctx) error {

	var app models.App
	if err := ctrl.db.Where("slug = ?", c.Params("app")).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var function models.CloudFunction
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("function"), app.ID).First(&function).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	return ctrl.executor.HandleRequest(function, c)
}
