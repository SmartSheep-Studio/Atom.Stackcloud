package main

import (
	"code.smartsheep.studio/atom/matrix/config"
	"code.smartsheep.studio/atom/matrix/datasource"
	"code.smartsheep.studio/atom/matrix/http"
	"code.smartsheep.studio/atom/matrix/logger"
	"code.smartsheep.studio/atom/matrix/services"
	"code.smartsheep.studio/atom/neutron/toolbox"
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
