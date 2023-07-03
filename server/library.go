package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"repo.smartsheep.studio/atom/matrix/datasource/models"
	"repo.smartsheep.studio/atom/matrix/server/middleware"
	"repo.smartsheep.studio/atom/nucleus/toolbox"
	"repo.smartsheep.studio/atom/nucleus/utils"
)

type LibraryController struct {
	db   *gorm.DB
	conn *toolbox.ExternalServiceConnection
	auth middleware.AuthHandler
}

func NewLibraryController(db *gorm.DB, conn *toolbox.ExternalServiceConnection, auth middleware.AuthHandler) *LibraryController {
	return &LibraryController{db, conn, auth}
}

func (ctrl *LibraryController) Map(router *fiber.App) {
	router.Get("/api/library", ctrl.auth(true), ctrl.list)
	router.Get("/api/library/own", ctrl.auth(true), ctrl.doesOwn)
	router.Post("/api/library/add", ctrl.auth(true), ctrl.add)
}

func (ctrl *LibraryController) list(c *fiber.Ctx) error {
	u := c.Locals("matrix-id").(*models.MatrixProfile)

	var items []models.MatrixLibraryItem
	if err := ctrl.db.Where("profile_id = ?", u.ID).Find(&items).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *LibraryController) doesOwn(c *fiber.Ctx) error {
	u := c.Locals("matrix-id").(*models.MatrixProfile)
	target := c.Query("app")

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ?", target).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var libraryCount int64
	if err := ctrl.db.Model(&models.MatrixLibraryItem{}).Where("profile_id = ? AND app_id = ?", u.ID, app.ID).Count(&libraryCount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return utils.ParseDataSourceError(err)
	} else if libraryCount <= 0 {
		return c.SendStatus(fiber.StatusOK)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func (ctrl *LibraryController) add(c *fiber.Ctx) error {
	u := c.Locals("matrix-id").(*models.MatrixProfile)

	var req struct {
		App string `json:"app" validate:"required"`
	}

	if err := utils.ParseRequestBody(c, &req); err != nil {
		return err
	}

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ?", req.App).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var libraryCount int64
	if err := ctrl.db.Model(&models.MatrixLibraryItem{}).Where("profile_id = ? AND app_id = ?", u.ID, app.ID).Count(&libraryCount).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ParseDataSourceError(err)
		}
	} else if libraryCount > 0 {
		return fiber.NewError(fiber.StatusForbidden, "already in the library")
	}

	item := models.MatrixLibraryItem{
		ProfileID: u.ID,
		AppID:     app.ID,
	}

	if err := ctrl.db.Save(&item).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(item)
	}
}
