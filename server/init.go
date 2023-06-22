package server

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"repo.smartsheep.studio/atom/matrix/renderer"
	"repo.smartsheep.studio/atom/matrix/server/middleware"
	"repo.smartsheep.studio/atom/nucleus/toolbox"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	flog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var conn *toolbox.ExternalServiceConnection
var server *fiber.App

func NewHttpServer(cycle fx.Lifecycle, cors middleware.CorsHandler, c *toolbox.ExternalServiceConnection) *fiber.App {
	conn = c

	// Create app
	server = fiber.New(fiber.Config{
		Prefork:               viper.GetBool("http.prefork"),
		CaseSensitive:         false,
		StrictRouting:         false,
		DisableStartupMessage: true,
		ServerHeader:          "Nucleus",
		AppName:               "Nucleus v2.0",
		BodyLimit:             viper.GetInt("http.max_body_size"),
	})

	// Apply global middleware
	server.Use(flog.New(flog.Config{
		Format: "${status} - ${latency} ${method} ${path}\n",
		Output: log.Logger,
	}))
	server.Use(cors())

	cycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := fmt.Sprintf("%s%d", viper.GetString("http.listen_address"), viper.GetInt("http.listen_port"))
			go func() {
				err := server.Listen(address)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to start http server")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown()
		},
	})

	return server
}

type HttpMap struct {
	Count uint32
}

func NewHttpMap(controllers []HttpController) *HttpMap {
	for _, controller := range controllers {
		controller.Map(server)
	}

	// Fallback not found api to nucleus
	server.All("/api/*", func(c *fiber.Ctx) error {
		uri := fmt.Sprintf("%s?%s", c.Request().URI().Path(), c.Request().URI().QueryArgs().String())
		return c.Redirect(conn.GetEndpointPath(uri), fiber.StatusFound)
	})

	// Serve static files
	server.Use("/", filesystem.New(filesystem.Config{
		Root:         renderer.GetHttpFS(),
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	return &HttpMap{Count: server.HandlersCount()}
}
