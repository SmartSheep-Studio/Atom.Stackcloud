package hypertext

import (
	"code.smartsheep.studio/atom/stackcloud/pkg/server/hypertext/controllers"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/hypertext/middleware"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("hypertext",
		middleware.Module(),
		controllers.Module(),

		fx.Invoke(fx.Annotate(MapControllers, fx.ParamTags(`group:"hypertext_controllers"`))),
		fx.Provide(NewHttpServer),
	)
}
