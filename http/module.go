package http

import (
	"code.smartsheep.studio/atom/matrix/http/controllers"
	"code.smartsheep.studio/atom/matrix/http/middleware"
	ctx "code.smartsheep.studio/atom/neutron/http/context"
	"go.uber.org/fx"
)

type HttpController interface {
	Map(router *ctx.App)
}

func AsController(f any) any {
	return fx.Annotate(f, fx.As(new(HttpController)), fx.ResultTags(`group:"controllers"`))
}

func Module() fx.Option {
	return fx.Module("http",
		middleware.Module(),
		controllers.Module(),
		
		fx.Invoke(fx.Annotate(MapControllers, fx.ParamTags(`group:"controllers"`))),
		fx.Provide(NewHttpServer),
	)
}
