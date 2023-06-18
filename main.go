package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"repo.smartsheep.studio/atom/lineupmarketplace/config"
	"repo.smartsheep.studio/atom/lineupmarketplace/datasource"
	"repo.smartsheep.studio/atom/lineupmarketplace/datasource/models"
	"repo.smartsheep.studio/atom/lineupmarketplace/logger"
	"repo.smartsheep.studio/atom/lineupmarketplace/server"
	"repo.smartsheep.studio/atom/lineupmarketplace/services"
	"repo.smartsheep.studio/atom/nucleus/toolbox"
)

func main() {
	fx.New(
		fx.WithLogger(logger.NewEventLogger),
		fx.Invoke(models.Migrate),

		datasource.Module(),
		services.Module(),
		server.Module(),

		fx.Provide(config.NewEndpointConnection),

		fx.Invoke(func(server *fiber.App, routes *server.HttpMap, endpoint *toolbox.ExternalServiceConnection) {
			log.Info().Msgf("You can open %s in your browser now", viper.GetString("general.base_url"))
		}),
	).Run()
}
