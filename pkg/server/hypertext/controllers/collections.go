package controllers

import (
	"code.smartsheep.studio/atom/stackcloud/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/hypertext/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CollectionController struct {
	db         *gorm.DB
	gatekeeper *middleware.AuthMiddleware
}

func NewCollectionController(db *gorm.DB, gatekeeper *middleware.AuthMiddleware) *CollectionController {
	return &CollectionController{db, gatekeeper}
}

func (ctrl *CollectionController) Map(router *fiber.App) {
	router.Get(
		"/api/apps/:app/records",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:records.collections"), hyperutils.GenPerms("records.collections.read")),
		ctrl.list,
	)
	router.Get(
		"/api/apps/:app/records/:collection",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:records.collections"), hyperutils.GenPerms("records.collections.read")),
		ctrl.get,
	)
	router.Post(
		"/api/apps/:app/records",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("create:records.collections"), hyperutils.GenPerms("records.collections.create")),
		ctrl.create,
	)
	router.Put(
		"/api/apps/:app/records/:collection",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("update:records.collections"), hyperutils.GenPerms("records.collections.update")),
		ctrl.update,
	)
	router.Delete(
		"/api/apps/:app/records/:collection",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("delete:records.collections"), hyperutils.GenPerms("records.collections.delete")),
		ctrl.delete,
	)
}

func (ctrl *CollectionController) list(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var collections []models.RecordCollection
	if err := ctrl.db.Where("app_id = ?", app.ID).Order("created_at desc").Find(&collections).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(collections)
	}
}

func (ctrl *CollectionController) get(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var collection models.RecordCollection
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("collection"), app.ID).First(&collection).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(collection)
	}
}

func (ctrl *CollectionController) create(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
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
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(collection)
	}
}

func (ctrl *CollectionController) update(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	var collection models.RecordCollection
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("collection"), app.ID).First(&collection).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	collection.Slug = req.Slug
	collection.Name = req.Name
	collection.Description = req.Description
	collection.Tags = datatypes.NewJSONSlice(req.Tags)

	if err := ctrl.db.Save(&collection).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(collection)
	}
}

func (ctrl *CollectionController) delete(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var collection models.RecordCollection
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("collection"), app.ID).First(&collection).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	if err := ctrl.db.Delete(&collection).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
