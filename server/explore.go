package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"repo.smartsheep.studio/atom/matrix/datasource/models"
	"repo.smartsheep.studio/atom/matrix/services"
	"repo.smartsheep.studio/atom/nucleus/utils"
)

type ExploreController struct {
	db      *gorm.DB
	service *services.ExploreService
}

func NewExploreController(db *gorm.DB, service *services.ExploreService) *ExploreController {
	return &ExploreController{db, service}
}

func (ctrl *ExploreController) Map(router *fiber.App) {
	router.Get("/api/explore/apps", ctrl.apps)
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

func (ctrl *ExploreController) posts(c *fiber.Ctx) error {
	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ?", c.Params("app")).First(&app).Error; err != nil {
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
	if err := ctrl.db.Where("slug = ?", c.Params("app")).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	items, err := ctrl.service.ExploreReleases(app.ID)
	if err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(items)
	}
}
