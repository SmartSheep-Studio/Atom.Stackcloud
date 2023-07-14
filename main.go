package main

import (
	"code.smartsheep.studio/atom/neutron/toolbox"
	"code.smartsheep.studio/atom/stackcloud/config"
	"code.smartsheep.studio/atom/stackcloud/datasource"
	"code.smartsheep.studio/atom/stackcloud/http"
	"code.smartsheep.studio/atom/stackcloud/logger"
	"code.smartsheep.studio/atom/stackcloud/services"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		logger.Module(),
		fx.WithLogger(logger.NewEventLogger),

		config.Module(),
		datasource.Module(),
		services.Module(),
		http.Module(),

		fx.Invoke(func(conf *viper.Viper, endpoint *toolbox.ExternalServiceConnection) {
			log.Info().Msgf("Your matrix instance is ready on: %s", conf.GetString("base_url"))
		}),
	).Run()
}
