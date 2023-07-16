package services

import (
	"encoding/json"
	"fmt"
	"time"

	tmodels "code.smartsheep.studio/atom/neutron/datasource/models"
	"code.smartsheep.studio/atom/neutron/http/context"
	"code.smartsheep.studio/atom/stackcloud/datasource/models"
	"github.com/dop251/goja"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FunctionService struct {
	db *gorm.DB
}

func NewFunctionService(db *gorm.DB) *FunctionService {
	return &FunctionService{db}
}

func (v *FunctionService) prepareRequestContext(vm *goja.Runtime, ctx *context.Ctx) error {
	var body map[string]any
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("request only accepts valid JSON format payload: %q", err))
	}

	request := vm.NewObject()
	request.Set("ip", ctx.IP())
	request.Set("body", body)
	request.Set("headers", ctx.GetReqHeaders())

	if ctx.Locals("principal-ok").(bool) {
		request.Set("principal", ctx.Locals("principal").(tmodels.User))
	} else {
		request.Set("principal", nil)
	}

	response := vm.NewObject()
	response.Set("text", func(call goja.FunctionCall) goja.Value {
		if err := ctx.Status(int(call.Argument(0).ToInteger())).SendString(call.Argument(1).String()); err != nil {
			panic(err)
		}
		return goja.Undefined()
	})
	response.Set("json", func(call goja.FunctionCall) goja.Value {
		if err := ctx.Status(int(call.Argument(0).ToInteger())).JSON(call.Argument(1).Export()); err != nil {
			panic(err)
		}
		return goja.Undefined()
	})

	vm.Set("request", request)
	vm.Set("response", response)

	return nil
}

func (v *FunctionService) prepareCollectionContext(vm *goja.Runtime, function models.CloudFunction, ctx *context.Ctx) error {
	var app models.App
	if err := v.db.Where("id = ?", function.AppID).First(&app).Error; err != nil {
		return ctx.DbError(err)
	}

	records := vm.NewObject()
	records.Set("find", func(call goja.FunctionCall) goja.Value {
		var collection models.RecordCollection
		if err := v.db.Where("slug = ? AND app_id = ?", call.Argument(0).String(), app.ID).First(&collection).Error; err != nil {
			return vm.NewGoError(err)
		}

		var record models.Record
		if err := v.db.Where("id = ?", call.Argument(1).ToInteger()).Error; err != nil {
			return vm.NewGoError(err)
		}

		return vm.ToValue(record)
	})
	records.Set("insert", func(call goja.FunctionCall) goja.Value {
		var collection models.RecordCollection
		if err := v.db.Where("slug = ? AND app_id = ?", call.Argument(0).String(), app.ID).First(&collection).Error; err != nil {
			return vm.NewGoError(err)
		}

		data, _ := json.Marshal(call.Argument(1).Export())
		record := models.Record{
			Payload:      datatypes.JSON(data),
			CollectionID: collection.ID,
		}
		if err := v.db.Save(&record).Error; err != nil {
			return vm.NewGoError(err)
		}

		return vm.ToValue(record)
	})
	records.Set("update", func(call goja.FunctionCall) goja.Value {
		var collection models.RecordCollection
		if err := v.db.Where("slug = ? AND app_id = ?", call.Argument(0).String(), app.ID).First(&collection).Error; err != nil {
			return vm.NewGoError(err)
		}

		var record models.Record
		if err := v.db.Where("id = ?", call.Argument(1).ToInteger()).Error; err != nil {
			return vm.NewGoError(err)
		}
		data, _ := json.Marshal(call.Argument(2).Export())
		record.Payload = datatypes.JSON(data)
		if err := v.db.Save(&record).Error; err != nil {
			return vm.NewGoError(err)
		}

		return vm.ToValue(record)
	})
	records.Set("delete", func(call goja.FunctionCall) goja.Value {
		var collection models.RecordCollection
		if err := v.db.Where("slug = ? AND app_id = ?", call.Argument(0).String(), app.ID).First(&collection).Error; err != nil {
			return vm.NewGoError(err)
		}

		var record models.Record
		if err := v.db.Where("id = ?", call.Argument(1).ToInteger()).Error; err != nil {
			return vm.NewGoError(err)
		}
		if err := v.db.Delete(&record).Error; err != nil {
			return vm.NewGoError(err)
		}

		return vm.ToValue(record)
	})

	vm.Set("records", records)

	return nil
}

func (v *FunctionService) HandleRequest(handler models.CloudFunction, ctx *context.Ctx) error {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	if err := v.prepareRequestContext(vm, ctx); err != nil {
		return err
	}
	if err := v.prepareCollectionContext(vm, handler, ctx); err != nil {
		return err
	}

	var err error
	time.AfterFunc(10*time.Second, func() {
		err = fiber.NewError(fiber.StatusInternalServerError, "function time limit exceeded 10s")
		vm.Interrupt("timeout")
	})

	if _, err := vm.RunString(handler.Script); err != nil {
		if ferr, ok := err.(*goja.Exception); ok {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("function internal error: %q", ferr))
		}
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("runtime internal error: %q", err))
	}

	return err
}
