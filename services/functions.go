package services

import (
	"fmt"
	"time"

	"code.smartsheep.studio/atom/neutron/http/context"
	"code.smartsheep.studio/atom/stackcloud/datasource/models"
	"github.com/dop251/goja"
	"github.com/gofiber/fiber/v2"
)

type FunctionService struct {
}

func NewFunctionService() *FunctionService {
	return &FunctionService{}
}

func (v *FunctionService) HandleRequest(handler models.CloudFunction, ctx *context.Ctx) error {
	vm := goja.New()

	var body map[string]any
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("request only accepts valid JSON format payload: %q", err))
	}

	request := vm.NewObject()
	request.Set("ip", ctx.IP())
	request.Set("body", body)
	request.Set("headers", ctx.GetReqHeaders())

	var err error
	response := vm.NewObject()
	response.Set("text", func(call goja.FunctionCall) goja.Value {
		err = ctx.Status(int(call.Argument(0).ToInteger())).SendString(call.Argument(1).String())
		return goja.Undefined()
	})
	response.Set("json", func(call goja.FunctionCall) goja.Value {
		err = ctx.Status(int(call.Argument(0).ToInteger())).JSON(call.Argument(1).Export())
		fmt.Println("responded!")
		return goja.Undefined()
	})

	vm.Set("request", request)
	vm.Set("response", response)

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
