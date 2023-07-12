package controllers

import (
	ctx "code.smartsheep.studio/atom/neutron/http/context"
	"errors"
	"github.com/gofiber/fiber/v2"

	"code.smartsheep.studio/atom/matrix/datasource/models"
	"code.smartsheep.studio/atom/matrix/http/middleware"
	"code.smartsheep.studio/atom/neutron/toolbox"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LibraryController struct {
	db   *gorm.DB
	conn *toolbox.ExternalServiceConnection
	auth middleware.AuthHandler
}

func NewLibraryController(db *gorm.DB, conn *toolbox.ExternalServiceConnection, auth middleware.AuthHandler) *LibraryController {
	return &LibraryController{db, conn, auth}
}

func (ctrl *LibraryController) Map(router *ctx.App) {
	router.Get("/api/library", ctrl.auth(true), ctrl.list)
	router.Get("/api/library/own", ctrl.auth(true), ctrl.doesOwn)
	router.Post("/api/library/add", ctrl.auth(true), ctrl.add)
}

func (ctrl *LibraryController) list(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var items []models.MatrixLibraryItem
	if err := ctrl.db.Where("account_id = ?", u.ID).Find(&items).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.JSON(items)
	}
}

func (ctrl *LibraryController) doesOwn(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)
	target := c.P.Query("app")

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ?", target).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var libraryCount int64
	if err := ctrl.db.Model(&models.MatrixLibraryItem{}).Where("account_id = ? AND app_id = ?", u.ID, app.ID).Count(&libraryCount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.P.SendStatus(fiber.StatusNoContent)
		}
		return c.DbError(err)
	} else if libraryCount <= 0 {
		return c.P.SendStatus(fiber.StatusOK)
	} else {
		return c.P.SendStatus(fiber.StatusNoContent)
	}
}

func (ctrl *LibraryController) add(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var req struct {
		App string `json:"app" validate:"required"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ?", req.App).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var libraryCount int64
	if err := ctrl.db.Model(&models.MatrixLibraryItem{}).Where("account_id = ? AND app_id = ?", u.ID, app.ID).Count(&libraryCount).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.DbError(err)
		}
	} else if libraryCount > 0 {
		return fiber.NewError(fiber.StatusForbidden, "already in the library")
	}

	item := models.MatrixLibraryItem{
		AccountID: u.ID,
		AppID:     app.ID,
		CloudSave: models.MatrixCloudSave{
			Name:    u.Nickname,
			Payload: datatypes.JSON([]byte("{}")),
		},
	}

	if err := ctrl.db.Save(&item).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.JSON(item)
	}
}
