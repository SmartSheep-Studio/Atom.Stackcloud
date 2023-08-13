package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type HypertextController interface {
	Map(router *fiber.App)
}

func AsController(f any) any {
	return fx.Annotate(f, fx.As(new(HypertextController)), fx.ResultTags(`group:"hypertext_controllers"`))
}

func Module() fx.Option {
	return fx.Module("hypertext.controllers",
		fx.Provide(
			AsController(NewStatusController),
			AsController(NewAppController),
			AsController(NewCollectionController),
			AsController(NewRecordController),
			AsController(NewFunctionController),
		),
	)
}
