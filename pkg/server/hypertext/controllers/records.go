package controllers

import (
	"code.smartsheep.studio/atom/stackcloud/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/hypertext/middleware"
	"encoding/json"

	"code.smartsheep.studio/atom/neutron/http/context"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RecordController struct {
	db         *gorm.DB
	gatekeeper *middleware.AuthMiddleware
}

func NewRecordController(db *gorm.DB, gatekeeper *middleware.AuthMiddleware) *RecordController {
	return &RecordController{db, gatekeeper}
}

func (ctrl *RecordController) Map(router *context.App) {
	router.Get(
		"/api/apps/:app/records/:collection/data",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:records.data"), hyperutils.GenPerms("records.data.read")),
		ctrl.list,
	)
	router.Get(
		"/api/apps/:app/records/:collection/data/:record",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:records.data"), hyperutils.GenPerms("records.data.read")),
		ctrl.get,
	)
	router.Post(
		"/api/apps/:app/records/:collection/data",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("create:records.data"), hyperutils.GenPerms("records.data.create")),
		ctrl.create,
	)
	router.Put(
		"/api/apps/:app/records/:collection/data/:record",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("update:records.data"), hyperutils.GenPerms("records.data.update")),
		ctrl.update,
	)
	router.Delete(
		"/api/apps/:app/records/:collection/data/:record",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("delete:records.data"), hyperutils.GenPerms("records.data.delete")),
		ctrl.delete,
	)
}

func (ctrl *RecordController) list(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var collection models.RecordCollection
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("collection"), app.ID).First(&collection).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var records []models.Record
	if err := ctrl.db.Where("collection_id = ?", collection.ID).Order("created_at desc").Find(&records).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(records)
	}
}

func (ctrl *RecordController) get(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var collection models.RecordCollection
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("collection"), app.ID).First(&collection).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var record models.Record
	if err := ctrl.db.Where("collection_id = ? AND id = ?", collection.ID, c.Params("record")).First(&record).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(record)
	}
}

func (ctrl *RecordController) create(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var collection models.RecordCollection
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("collection"), app.ID).First(&collection).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var req struct {
		Payload map[string]any `json:"payload"`
	}

	var payload []byte
	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	} else {
		payload, _ = json.Marshal(req.Payload)
	}

	record := models.Record{
		Payload:      datatypes.JSON(payload),
		CollectionID: collection.ID,
	}

	if err := ctrl.db.Save(&record).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(record)
	}
}

func (ctrl *RecordController) update(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var collection models.RecordCollection
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("collection"), app.ID).First(&collection).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var req struct {
		Payload map[string]any `json:"payload"`
	}

	var payload []byte
	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	} else {
		payload, _ = json.Marshal(req.Payload)
	}

	var record models.Record
	if err := ctrl.db.Where("collection_id = ? AND id = ?", collection.ID, c.Params("record")).First(&record).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	record.Payload = datatypes.JSON(payload)

	if err := ctrl.db.Save(&record).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(record)
	}
}

func (ctrl *RecordController) delete(c *fiber.Ctx) error {

	u := c.Locals("stackcloud-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var collection models.RecordCollection
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("collection"), app.ID).First(&collection).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var record models.Record
	if err := ctrl.db.Where("collection_id = ? AND id = ?", collection.ID, c.Params("record")).First(&record).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	if err := ctrl.db.Delete(&record).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
