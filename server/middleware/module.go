package middleware

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("server.middleware",
		fx.Provide(NewAuth),
		fx.Provide(NewCors),
	)
}
