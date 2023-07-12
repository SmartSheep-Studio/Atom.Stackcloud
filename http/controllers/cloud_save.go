package controllers

import (
	ctx "code.smartsheep.studio/atom/neutron/http/context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"

	"code.smartsheep.studio/atom/matrix/datasource/models"
	"code.smartsheep.studio/atom/matrix/http/middleware"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CloudSaveController struct {
	db   *gorm.DB
	auth middleware.AuthHandler
}

func NewCloudSaveController(db *gorm.DB, auth middleware.AuthHandler) *CloudSaveController {
	return &CloudSaveController{db, auth}
}

func (ctrl *CloudSaveController) Map(router *ctx.App) {
	router.Get("/api/apps/:app/cloud-save", ctrl.auth(true), ctrl.get)
	router.Put("/api/apps/:app/cloud-save", ctrl.auth(true), ctrl.update)
	router.Put("/api/apps/:app/cloud-save/name", ctrl.auth(true), ctrl.updateInfo)
}

func (ctrl *CloudSaveController) get(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ?", c.P.Params("app")).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var library models.MatrixLibraryItem
	if err := ctrl.db.Where("app_id = ? AND user_id = ?", app.ID, u.ID).Preload("CloudSave").First(&library).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(fiber.StatusForbidden, "you haven't that app")
		} else {
			return c.DbError(err)
		}
	}

	return c.P.JSON(library.CloudSave)
}

func (ctrl *CloudSaveController) update(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var req map[string]any
	if err := c.BindBody(&req); err != nil {
		return err
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("you need provide a valid json format payload: %q", err))
	}

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ?", c.P.Params("app")).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var library models.MatrixLibraryItem
	if err := ctrl.db.Where("app_id = ? AND user_id = ?", app.ID, u.ID).Preload("CloudSave").First(&library).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(fiber.StatusForbidden, "you haven't that app")
		} else {
			return c.DbError(err)
		}
	}

	library.CloudSave.Payload = datatypes.JSON(data)

	if err := ctrl.db.Save(&library.CloudSave).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.JSON(library.CloudSave)
	}
}

func (ctrl *CloudSaveController) updateInfo(c *ctx.Ctx) error {
	u := c.P.Locals("matrix-id").(*models.MatrixAccount)

	var req struct {
		Name string `json:"name" validate:"required"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ?", c.P.Params("app")).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var library models.MatrixLibraryItem
	if err := ctrl.db.Where("app_id = ? AND user_id = ?", app.ID, u.ID).Preload("CloudSave").First(&library).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(fiber.StatusForbidden, "you haven't that app")
		} else {
			return c.DbError(err)
		}
	}

	library.CloudSave.Name = req.Name

	if err := ctrl.db.Save(&library.CloudSave).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.P.JSON(library.CloudSave)
	}
}
