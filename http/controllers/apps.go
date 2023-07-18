package controllers

import (
	"code.smartsheep.studio/atom/neutron/http/context"
	"code.smartsheep.studio/atom/stackcloud/datasource/models"
	"code.smartsheep.studio/atom/stackcloud/http/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AppController struct {
	db   *gorm.DB
	auth middleware.AuthHandler
}

func NewAppController(db *gorm.DB, auth middleware.AuthHandler) *AppController {
	return &AppController{db, auth}
}

func (ctrl *AppController) Map(router *context.App) {
	router.Get("/api/apps", ctrl.auth(true, "stackcloud.apps.read", "stackcloud.apps.read"), ctrl.list)
	router.Get("/api/apps/:app", ctrl.auth(true, "stackcloud.apps.read", "stackcloud.apps.read"), ctrl.get)
	router.Post("/api/apps", ctrl.auth(true, "stackcloud.apps.create", "stackcloud.apps.create"), ctrl.create)
	router.Put("/api/apps/:app", ctrl.auth(true, "stackcloud.apps.update", "stackcloud.apps.update"), ctrl.update)
	router.Delete("/api/apps/:app", ctrl.auth(true, "stackcloud.apps.delete", "stackcloud.apps.delete"), ctrl.delete)
}

func (ctrl *AppController) list(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var apps []models.App
	if err := ctrl.db.Where("account_id = ?", u.ID).Find(&apps).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(apps)
	}
}

func (ctrl *AppController) get(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(app)
	}
}

func (ctrl *AppController) create(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	app := models.App{
		Slug:        req.Slug,
		Tags:        datatypes.NewJSONSlice(req.Tags),
		Name:        req.Name,
		Description: req.Description,
		AccountID:   u.ID,
	}

	if err := ctrl.db.Save(&app).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(app)
	}
}

func (ctrl *AppController) update(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	app.Tags = datatypes.NewJSONSlice(req.Tags)
	app.Slug = req.Slug
	app.Name = req.Name
	app.Description = req.Description

	if err := ctrl.db.Save(&app).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(app)
	}
}

func (ctrl *AppController) delete(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	if err := ctrl.db.Delete(&app).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
