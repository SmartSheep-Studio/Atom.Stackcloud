package controllers

import (
	"code.smartsheep.studio/atom/matrix/datasource/models"
	"code.smartsheep.studio/atom/matrix/http/middleware"
	ctx "code.smartsheep.studio/atom/neutron/http/context"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PostController struct {
	db   *gorm.DB
	auth middleware.AuthHandler
}

func NewPostController(db *gorm.DB, auth middleware.AuthHandler) *PostController {
	return &PostController{db, auth}
}

func (ctrl *PostController) Map(router *ctx.App) {
	router.Get("/api/apps/:app/posts", ctrl.auth(true), ctrl.list)
	router.Get("/api/apps/:app/posts/:post", ctrl.auth(true), ctrl.get)
	router.Post("/api/apps/:app/posts", ctrl.auth(true), ctrl.create)
	router.Put("/api/apps/:app/posts/:post", ctrl.auth(true), ctrl.update)
	router.Delete("/api/apps/:app/posts/:post", ctrl.auth(true), ctrl.delete)
}

func (ctrl *PostController) list(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.P.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var posts []models.MatrixPost
	if err := ctrl.db.Where("app_id = ?", app.ID).Order("created_at desc").Find(&posts).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.JSON(posts)
	}
}

func (ctrl *PostController) get(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.P.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var post models.MatrixPost
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.P.Params("post"), app.ID).First(&post).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.JSON(post)
	}
}

func (ctrl *PostController) create(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.P.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Title       string   `json:"title" validate:"required"`
		Type        string   `json:"type" validate:"required"`
		Content     string   `json:"content"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := c.BindBody(&req); err != nil {
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
		return c.DbError(err)
	} else {
		return c.P.JSON(post)
	}
}

func (ctrl *PostController) update(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.P.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var req struct {
		Title       string   `json:"title" validate:"required"`
		Type        string   `json:"type" validate:"required"`
		Content     string   `json:"content"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	var post models.MatrixPost
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.P.Params("post"), app.ID).First(&post).Error; err != nil {
		return c.DbError(err)
	}

	post.Title = req.Title
	post.Type = req.Type
	post.Content = req.Content
	post.Tags = datatypes.NewJSONSlice(req.Tags)
	post.IsPublished = req.IsPublished

	if err := ctrl.db.Save(&post).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.JSON(post)
	}
}

func (ctrl *PostController) delete(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.P.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var post models.MatrixPost
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.P.Params("post"), app.ID).First(&post).Error; err != nil {
		return c.DbError(err)
	}

	if err := ctrl.db.Delete(&post).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.SendStatus(fiber.StatusNoContent)
	}
}
