package middleware

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("http.middleware",
		fx.Provide(NewAuth),
		fx.Provide(NewCors),
	)
}
