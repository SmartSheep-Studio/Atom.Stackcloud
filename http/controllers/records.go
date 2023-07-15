package controllers

import (
	"encoding/json"

	"code.smartsheep.studio/atom/neutron/http/context"
	"code.smartsheep.studio/atom/stackcloud/datasource/models"
	"code.smartsheep.studio/atom/stackcloud/http/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RecordController struct {
	db   *gorm.DB
	auth middleware.AuthHandler
}

func NewRecordController(db *gorm.DB, auth middleware.AuthHandler) *RecordController {
	return &RecordController{db, auth}
}

func (ctrl *RecordController) Map(router *context.App) {
	router.Get("/api/apps/:app/records/:collection/data", ctrl.auth(true, "records.data.read", "stackcloud.records.data.read"), ctrl.list)
	router.Get("/api/apps/:app/records/:collection/data/:record", ctrl.auth(true, "records.data.read", "stackcloud.records.data.read"), ctrl.get)
	router.Post("/api/apps/:app/records/:collection/data", ctrl.auth(true, "records.data.create", "stackcloud.records.data.create"), ctrl.create)
	router.Put("/api/apps/:app/records/:collection/data/:record", ctrl.auth(true, "records.data.update", "stackcloud.records.data.update"), ctrl.update)
	router.Delete("/api/apps/:app/records/:collection/data/:record", ctrl.auth(true, "records.data.delete", "stackcloud.records.data.delete"), ctrl.delete)
}

func (ctrl *RecordController) list(ctx *fiber.Ctx) error {
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

	var records []models.Record
	if err := ctrl.db.Where("collection_id = ?", collection.ID).Order("created_at desc").Find(&records).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(records)
	}
}

func (ctrl *RecordController) get(ctx *fiber.Ctx) error {
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

	var record models.Record
	if err := ctrl.db.Where("collection_id = ? AND id = ?", collection.ID, c.Params("record")).First(&record).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(record)
	}
}

func (ctrl *RecordController) create(ctx *fiber.Ctx) error {
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

	var req struct {
		Payload map[string]any `json:"payload"`
	}

	var payload []byte
	if err := c.BindBody(&req); err != nil {
		return err
	} else {
		payload, _ = json.Marshal(req.Payload)
	}

	record := models.Record{
		Payload:      datatypes.JSON(payload),
		CollectionID: collection.ID,
	}

	if err := ctrl.db.Save(&record).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(record)
	}
}

func (ctrl *RecordController) update(ctx *fiber.Ctx) error {
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

	var req struct {
		Payload map[string]any `json:"payload"`
	}

	var payload []byte
	if err := c.BindBody(&req); err != nil {
		return err
	} else {
		payload, _ = json.Marshal(req.Payload)
	}

	var record models.Record
	if err := ctrl.db.Where("collection_id = ? AND id = ?", collection.ID, c.Params("record")).First(&record).Error; err != nil {
		return c.DbError(err)
	}

	record.Payload = datatypes.JSON(payload)

	if err := ctrl.db.Save(&record).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(record)
	}
}

func (ctrl *RecordController) delete(ctx *fiber.Ctx) error {
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

	var record models.Record
	if err := ctrl.db.Where("collection_id = ? AND id = ?", collection.ID, c.Params("record")).First(&record).Error; err != nil {
		return c.DbError(err)
	}

	if err := ctrl.db.Delete(&record).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
