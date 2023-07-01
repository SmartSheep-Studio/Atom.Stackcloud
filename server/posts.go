package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"repo.smartsheep.studio/atom/matrix/datasource/models"
	"repo.smartsheep.studio/atom/matrix/server/middleware"
	"repo.smartsheep.studio/atom/nucleus/utils"
)

type PostController struct {
	db   *gorm.DB
	auth middleware.AuthHandler
}

func NewPostController(db *gorm.DB, auth middleware.AuthHandler) *PostController {
	return &PostController{db, auth}
}

func (ctrl *PostController) Map(router *fiber.App) {
	router.Get("/api/apps/:app/posts", ctrl.auth(true), ctrl.list)
	router.Get("/api/apps/:app/posts/:post", ctrl.auth(true), ctrl.get)
	router.Post("/api/apps/:app/posts", ctrl.auth(true), ctrl.create)
	router.Put("/api/apps/:app/posts/:post", ctrl.auth(true), ctrl.update)
	router.Delete("/api/apps/:app/posts/:post", ctrl.auth(true), ctrl.delete)
}

func (ctrl *PostController) list(c *fiber.Ctx) error {
	u := c.Locals("matrix-prof").(*models.MatrixProfile)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND profile_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var posts []models.MatrixPost
	if err := ctrl.db.Where("app_id = ?", app.ID).Order("created_at desc").Find(&posts).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(posts)
	}
}

func (ctrl *PostController) get(c *fiber.Ctx) error {
	u := c.Locals("matrix-prof").(*models.MatrixProfile)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND profile_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var post models.MatrixPost
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("post"), app.ID).First(&post).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(post)
	}
}

func (ctrl *PostController) create(c *fiber.Ctx) error {
	u := c.Locals("matrix-prof").(*models.MatrixProfile)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND profile_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Title       string   `json:"title" validate:"required"`
		Type        string   `json:"type" validate:"required"`
		Content     string   `json:"content"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := utils.ParseRequestBody(c, &req); err != nil {
		return err
	}

	post := models.MatrixPost{
		Slug:        req.Slug,
		Type:        req.Type,
		Title:       req.Title,
		Content:     req.Content,
		Tags:        datatypes.NewJSONSlice(req.Tags),
		AppID:       app.ID,
		IsPublished: req.IsPublished,
	}

	if err := ctrl.db.Save(&post).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(post)
	}
}

func (ctrl *PostController) update(c *fiber.Ctx) error {
	u := c.Locals("matrix-prof").(*models.MatrixProfile)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND profile_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var req struct {
		Title       string   `json:"title" validate:"required"`
		Type        string   `json:"type" validate:"required"`
		Content     string   `json:"content"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := utils.ParseRequestBody(c, &req); err != nil {
		return err
	}

	var post models.MatrixPost
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("post"), app.ID).First(&post).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	post.Title = req.Title
	post.Type = req.Type
	post.Content = req.Content
	post.Tags = datatypes.NewJSONSlice(req.Tags)
	post.IsPublished = req.IsPublished

	if err := ctrl.db.Save(&post).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(post)
	}
}

func (ctrl *PostController) delete(c *fiber.Ctx) error {
	u := c.Locals("matrix-prof").(*models.MatrixProfile)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND profile_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var post models.MatrixPost
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("post"), app.ID).First(&post).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	if err := ctrl.db.Delete(&post).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
