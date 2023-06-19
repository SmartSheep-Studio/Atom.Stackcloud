package server

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"repo.smartsheep.studio/atom/matrix/server/middleware"
)

type HttpController interface {
	Map(router *fiber.App)
}

func AsController(f any) any {
	return fx.Annotate(f, fx.As(new(HttpController)), fx.ResultTags(`group:"controllers"`))
}

func Module() fx.Option {
	return fx.Module("server",
		middleware.Module(),

		fx.Provide(NewHttpServer),
		fx.Provide(
			AsController(NewStatusController),
			AsController(NewAuthController),
			AsController(NewExploreController),
			AsController(NewAppController),

			fx.Annotate(NewHttpMap, fx.ParamTags(`group:"controllers"`)),
		),
	)
}
