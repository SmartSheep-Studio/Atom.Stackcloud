package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"repo.smartsheep.studio/atom/matrix/datasource/models"
	"repo.smartsheep.studio/atom/matrix/server/middleware"
	"repo.smartsheep.studio/atom/nucleus/utils"
)

type AppController struct {
	db   *gorm.DB
	auth middleware.AuthHandler
}

func NewAppController(db *gorm.DB, auth middleware.AuthHandler) *AppController {
	return &AppController{db, auth}
}

func (ctrl *AppController) Map(router *fiber.App) {
	router.Get("/api/apps", ctrl.auth(true), ctrl.list)
	router.Get("/api/apps/:app", ctrl.auth(true), ctrl.get)
	router.Post("/api/apps", ctrl.auth(true), ctrl.create)
	router.Put("/api/apps/:app", ctrl.auth(true), ctrl.update)
	router.Delete("/api/apps/:app", ctrl.auth(true), ctrl.delete)
}

func (ctrl *AppController) list(c *fiber.Ctx) error {
	u := c.Locals("matrix-id").(*models.MatrixAccount)

	var apps []models.MatrixApp
	if err := ctrl.db.Where("account_id = ?", u.ID).Find(&apps).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(apps)
	}
}

func (ctrl *AppController) get(c *fiber.Ctx) error {
	u := c.Locals("matrix-id").(*models.MatrixAccount)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(app)
	}
}

func (ctrl *AppController) create(c *fiber.Ctx) error {
	u := c.Locals("matrix-id").(*models.MatrixAccount)

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Details     string   `json:"details"`
		Url         string   `json:"url"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := utils.ParseRequestBody(c, &req); err != nil {
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
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(app)
	}
}

func (ctrl *AppController) update(c *fiber.Ctx) error {
	u := c.Locals("matrix-id").(*models.MatrixAccount)

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Details     string   `json:"details"`
		Url         string   `json:"url"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := utils.ParseRequestBody(c, &req); err != nil {
		return err
	}

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	app.Url = req.Url
	app.Tags = datatypes.NewJSONSlice(req.Tags)
	app.Slug = req.Slug
	app.Name = req.Name
	app.Description = req.Description
	app.Details = req.Details
	app.IsPublished = req.IsPublished

	if err := ctrl.db.Save(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(app)
	}
}

func (ctrl *AppController) delete(c *fiber.Ctx) error {
	u := c.Locals("matrix-id").(*models.MatrixAccount)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	if err := ctrl.db.Delete(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
