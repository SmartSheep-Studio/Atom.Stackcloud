package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"repo.smartsheep.studio/atom/matrix/datasource/models"
	"repo.smartsheep.studio/atom/matrix/services"
	"repo.smartsheep.studio/atom/nucleus/toolbox"
	"repo.smartsheep.studio/atom/nucleus/utils"
)

type ExploreController struct {
	db      *gorm.DB
	conn    *toolbox.ExternalServiceConnection
	service *services.ExploreService
}

func NewExploreController(db *gorm.DB, conn *toolbox.ExternalServiceConnection, service *services.ExploreService) *ExploreController {
	return &ExploreController{db, conn, service}
}

func (ctrl *ExploreController) Map(router *fiber.App) {
	router.Get("/api/explore/apps", ctrl.apps)
	router.Get("/api/explore/apps/:app", ctrl.app)
	router.Get("/api/explore/apps/:app/posts", ctrl.posts)
	router.Get("/api/explore/apps/:app/releases", ctrl.releases)
}

func (ctrl *ExploreController) apps(c *fiber.Ctx) error {
	items, err := ctrl.service.ExploreApps()
	if err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *ExploreController) app(c *fiber.Ctx) error {
	var app models.MatrixApp
	id, _ := c.ParamsInt("app", 0)
	if err := ctrl.db.Where("slug = ? OR id = ?", c.Params("app"), id).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	items, err := ctrl.service.ExploreApp(app.ID)
	if err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *ExploreController) posts(c *fiber.Ctx) error {
	var app models.MatrixApp
	id, _ := c.ParamsInt("app", 0)
	if err := ctrl.db.Where("slug = ? OR id = ?", c.Params("app"), id).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	items, err := ctrl.service.ExplorePosts(app.ID)
	if err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *ExploreController) releases(c *fiber.Ctx) error {
	var app models.MatrixApp
	id, _ := c.ParamsInt("app", 0)
	if err := ctrl.db.Where("slug = ? OR id = ?", c.Params("app"), id).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	items, err := ctrl.service.ExploreReleases(app.ID)
	if err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(items)
	}
}
