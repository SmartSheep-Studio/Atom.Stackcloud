package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"repo.smartsheep.studio/atom/matrix/datasource/models"
	"repo.smartsheep.studio/atom/matrix/server/middleware"
	"repo.smartsheep.studio/atom/nucleus/utils"
)

type ReleaseController struct {
	db   *gorm.DB
	auth middleware.AuthHandler
}

func NewReleaseController(db *gorm.DB, auth middleware.AuthHandler) *ReleaseController {
	return &ReleaseController{db, auth}
}

func (ctrl *ReleaseController) Map(router *fiber.App) {
	router.Get("/api/apps/:app/releases", ctrl.auth(true), ctrl.list)
	router.Get("/api/apps/:app/releases/:release", ctrl.auth(true), ctrl.get)
	router.Post("/api/apps/:app/releases", ctrl.auth(true), ctrl.create)
	router.Put("/api/apps/:app/releases/:release", ctrl.auth(true), ctrl.update)
	router.Delete("/api/apps/:app/releases/:release", ctrl.auth(true), ctrl.delete)
}

func (ctrl *ReleaseController) list(c *fiber.Ctx) error {
	u := c.Locals("matrix-prof").(*models.MatrixProfile)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND profile_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var releases []models.MatrixRelease
	if err := ctrl.db.Where("app_id = ?", app.ID).Order("created_at desc").Preload("Post").Find(&releases).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(releases)
	}
}

func (ctrl *ReleaseController) get(c *fiber.Ctx) error {
	u := c.Locals("matrix-prof").(*models.MatrixProfile)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND profile_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var shop models.MatrixApp
	if err := ctrl.db.Where("id = ? AND app_id = ?", c.Params("release"), app.ID).Preload("Post").First(&shop).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(shop)
	}
}

func (ctrl *ReleaseController) create(c *fiber.Ctx) error {
	u := c.Locals("matrix-prof").(*models.MatrixProfile)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND profile_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Type        string   `json:"type" validate:"required"`
		Description string   `json:"description"`
		Details     string   `json:"details"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := utils.ParseRequestBody(c, &req); err != nil {
		return err
	}

	release := models.MatrixRelease{
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		Post: models.MatrixPost{
			Type:        req.Type,
			Title:       req.Name,
			Content:     req.Details,
			Tags:        datatypes.NewJSONSlice(req.Tags),
			IsPublished: req.IsPublished,
			AppID:       app.ID,
		},
		IsPublished: req.IsPublished,
		AppID:       app.ID,
	}

	if err := ctrl.db.Save(&release).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(release)
	}
}

func (ctrl *ReleaseController) update(c *fiber.Ctx) error {
	u := c.Locals("matrix-prof").(*models.MatrixProfile)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND profile_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Type        string   `json:"type" validate:"required"`
		Description string   `json:"description"`
		Details     string   `json:"details"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := utils.ParseRequestBody(c, &req); err != nil {
		return err
	}

	var release models.MatrixRelease
	if err := ctrl.db.Where("id = ? AND app_id = ?", c.Params("release"), app.ID).Preload("Post").First(&release).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		ctrl.db.Unscoped().Delete(&release.Post)
	}

	release.Slug = req.Slug
	release.Name = req.Name
	release.Description = req.Description
	release.IsPublished = req.IsPublished
	release.Post = models.MatrixPost{
		Type:        req.Type,
		Title:       req.Name,
		Content:     req.Details,
		Tags:        datatypes.NewJSONSlice(req.Tags),
		IsPublished: req.IsPublished,
		AppID:       app.ID,
	}

	if err := ctrl.db.Save(&release).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(release)
	}
}

func (ctrl *ReleaseController) delete(c *fiber.Ctx) error {
	u := c.Locals("matrix-prof").(*models.MatrixProfile)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND profile_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var release models.MatrixRelease
	if err := ctrl.db.Where("id = ? AND app_id = ?", c.Params("release"), app.ID).Preload("Post").First(&release).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	if err := ctrl.db.Delete(&release.Post).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	if err := ctrl.db.Delete(&release).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
