package main

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/config"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/datasource"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/hypertext"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/logger"
	"code.smartsheep.studio/atom/stackcloud/pkg/server/services"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"go.uber.org/fx"
	"os"
)

func main() {
	instance := &cli.App{
		Name:    "Matrix Server",
		Usage:   "An Open-Source Apps and Games marketplace",
		Version: "1.0-SNAPSHOT",
		Commands: []cli.Command{
			{
				Name:  "serve",
				Usage: "Start server",
				Action: func(c *cli.Context) error {
					log.Info().Msgf("You are running Matrix %s!", "SNAPSHOT-1.0")
					fx.New(
						logger.Module(),
						fx.WithLogger(logger.NewEventLogger),

						config.Module(),
						datasource.Module(),
						services.Module(),
						hypertext.Module(),

						fx.Invoke(func(conf *viper.Viper, endpoint *subapps.HeLiCoPtErConnection) {
							log.Info().Msgf("Your matrix instance is ready on: %s", conf.GetString("base_url"))
						}),
					).Run()
					return nil
				},
			},
		},
	}

	if err := instance.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run bedrock server.")
	}
}
