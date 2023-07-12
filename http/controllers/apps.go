package controllers

import (
	"code.smartsheep.studio/atom/matrix/datasource/models"
	"code.smartsheep.studio/atom/matrix/http/middleware"
	ctx "code.smartsheep.studio/atom/neutron/http/context"
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

func (ctrl *AppController) Map(router *ctx.App) {
	router.Get("/api/apps", ctrl.auth(true), ctrl.list)
	router.Get("/api/apps/:app", ctrl.auth(true), ctrl.get)
	router.Post("/api/apps", ctrl.auth(true), ctrl.create)
	router.Put("/api/apps/:app", ctrl.auth(true), ctrl.update)
	router.Delete("/api/apps/:app", ctrl.auth(true), ctrl.delete)
}

func (ctrl *AppController) list(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var apps []models.MatrixApp
	if err := ctrl.db.Where("account_id = ?", u.ID).Find(&apps).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.JSON(apps)
	}
}

func (ctrl *AppController) get(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.P.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.JSON(app)
	}
}

func (ctrl *AppController) create(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Details     string   `json:"details"`
		Url         string   `json:"url"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	app := models.MatrixApp{
		Slug:        req.Slug,
		Url:         req.Url,
		Tags:        datatypes.NewJSONSlice(req.Tags),
		Name:        req.Name,
		Description: req.Description,
		Details:     req.Details,
		IsPublished: req.IsPublished,
		AccountID:   u.ID,
	}

	if err := ctrl.db.Save(&app).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.JSON(app)
	}
}

func (ctrl *AppController) update(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Details     string   `json:"details"`
		Url         string   `json:"url"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.P.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	app.Url = req.Url
	app.Tags = datatypes.NewJSONSlice(req.Tags)
	app.Slug = req.Slug
	app.Name = req.Name
	app.Description = req.Description
	app.Details = req.Details
	app.IsPublished = req.IsPublished

	if err := ctrl.db.Save(&app).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.JSON(app)
	}
}

func (ctrl *AppController) delete(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.P.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	if err := ctrl.db.Delete(&app).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.SendStatus(fiber.StatusNoContent)
	}
}
