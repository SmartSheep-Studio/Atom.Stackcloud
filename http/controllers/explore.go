package controllers

import (
	"code.smartsheep.studio/atom/matrix/datasource/models"
	"code.smartsheep.studio/atom/matrix/services"
	"code.smartsheep.studio/atom/neutron/http/context"
	"code.smartsheep.studio/atom/neutron/toolbox"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ExploreController struct {
	db      *gorm.DB
	conn    *toolbox.ExternalServiceConnection
	service *services.ExploreService
}

func NewExploreController(db *gorm.DB, conn *toolbox.ExternalServiceConnection, service *services.ExploreService) *ExploreController {
	return &ExploreController{db, conn, service}
}

func (ctrl *ExploreController) Map(router *context.App) {
	router.Get("/api/explore/apps", ctrl.apps)
	router.Get("/api/explore/apps/:app", ctrl.app)
	router.Get("/api/explore/apps/:app/posts", ctrl.posts)
	router.Get("/api/explore/apps/:app/releases", ctrl.releases)
}

func (ctrl *ExploreController) apps(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	items, err := ctrl.service.ExploreApps()
	if err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *ExploreController) app(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	var app models.App
	id, _ := c.ParamsInt("app", 0)
	if err := ctrl.db.Where("slug = ? OR id = ?", c.Params("app"), id).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	items, err := ctrl.service.ExploreApp(app.ID)
	if err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *ExploreController) posts(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	var app models.App
	id, _ := c.ParamsInt("app", 0)
	if err := ctrl.db.Where("slug = ? OR id = ?", c.Params("app"), id).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	items, err := ctrl.service.ExplorePosts(app.ID)
	if err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *ExploreController) releases(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	var app models.App
	id, _ := c.ParamsInt("app", 0)
	if err := ctrl.db.Where("slug = ? OR id = ?", c.Params("app"), id).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	items, err := ctrl.service.ExploreReleases(app.ID)
	if err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(items)
	}
}
