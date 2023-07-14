package http

import (
	"code.smartsheep.studio/atom/stackcloud/http/controllers"
	"code.smartsheep.studio/atom/stackcloud/http/middleware"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("http",
		middleware.Module(),
		controllers.Module(),

		fx.Invoke(fx.Annotate(MapControllers, fx.ParamTags(`group:"http"`))),
		fx.Provide(NewHttpServer),
	)
}
