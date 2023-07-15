package controllers

import (
	"code.smartsheep.studio/atom/neutron/http/context"
	"code.smartsheep.studio/atom/stackcloud/datasource/models"
	"code.smartsheep.studio/atom/stackcloud/http/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CollectionController struct {
	db   *gorm.DB
	auth middleware.AuthHandler
}

func NewCollectionController(db *gorm.DB, auth middleware.AuthHandler) *CollectionController {
	return &CollectionController{db, auth}
}

func (ctrl *CollectionController) Map(router *context.App) {
	router.Get("/api/apps/:app/records", ctrl.auth(true, "records.collection.read", "stackcloud.records.collection.read"), ctrl.list)
	router.Get("/api/apps/:app/records/:collection", ctrl.auth(true, "records.collection.read", "stackcloud.records.collection.read"), ctrl.get)
	router.Post("/api/apps/:app/records", ctrl.auth(true, "records.collection.create", "stackcloud.records.collection.create"), ctrl.create)
	router.Put("/api/apps/:app/records/:collection", ctrl.auth(true, "records.collection.update", "stackcloud.records.collection.update"), ctrl.update)
	router.Delete("/api/apps/:app/records/:collection", ctrl.auth(true, "records.collection.delete", "stackcloud.records.collection.delete"), ctrl.delete)
}

func (ctrl *CollectionController) list(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var collections []models.RecordCollection
	if err := ctrl.db.Where("app_id = ?", app.ID).Order("created_at desc").Find(&collections).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(collections)
	}
}

func (ctrl *CollectionController) get(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var collection models.RecordCollection
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("collection"), app.ID).First(&collection).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(collection)
	}
}

func (ctrl *CollectionController) create(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	collection := models.RecordCollection{
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		Tags:        datatypes.NewJSONSlice(req.Tags),
		AppID:       app.ID,
	}

	if err := ctrl.db.Save(&collection).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(collection)
	}
}

func (ctrl *CollectionController) update(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	var collection models.RecordCollection
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("collection"), app.ID).First(&collection).Error; err != nil {
		return c.DbError(err)
	}

	collection.Slug = req.Slug
	collection.Name = req.Name
	collection.Description = req.Description
	collection.Tags = datatypes.NewJSONSlice(req.Tags)

	if err := ctrl.db.Save(&collection).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(collection)
	}
}

func (ctrl *CollectionController) delete(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var collection models.RecordCollection
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("collection"), app.ID).First(&collection).Error; err != nil {
		return c.DbError(err)
	}

	if err := ctrl.db.Delete(&collection).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
