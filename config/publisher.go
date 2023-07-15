package config

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"

	tmodels "code.smartsheep.studio/atom/neutron/datasource/models"
	"code.smartsheep.studio/atom/neutron/toolbox"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func NewEndpointConnection(cycle fx.Lifecycle) *toolbox.ExternalServiceConnection {
	connection := &toolbox.ExternalServiceConnection{}

	cycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			conn, err := toolbox.PublishService(viper.GetString("endpoints"), viper.GetString("mount_key"), toolbox.ExternalServiceRequest{
				Name:        "Stackcloud",
				InstanceID:  viper.GetString("instance_id"),
				PackageID:   "code.smartsheep.studio/atom/stackcloud",
				Description: "A developer-friendly cloud computing and serverless services platform.",
				Address:     viper.GetString("base_url"),
				Options: tmodels.ExternalServiceOptions{
					Pages: []tmodels.ExternalPage{
						{
							To:      viper.GetString("base_url"),
							Title:   "Matrix",
							Name:    "matrix",
							Icon:    "mdi-store",
							Builtin: false,
							Visible: true,
							Meta: map[string]any{
								"gatekeeper": map[string]any{
									"must": true,
								},
							},
						},
					},
					Requirements: []string{"oauth"},
					Properties: fiber.Map{
						"oauth.urls":      []string{viper.GetString("base_url")},
						"oauth.callbacks": []string{fmt.Sprintf("%s/api/auth/callback", viper.GetString("base_url"))},
					},
				},
			})

			if err != nil {
				return err
			} else {
				connection.Configuration = conn.Configuration
				connection.Service = conn.Service
				connection.Additional = conn.Additional

				log.Info().Fields(connection.Service).Msg("Successfully published service into endpoints!")
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return toolbox.DepublishService(viper.GetString("general.endpoints"), viper.GetString("general.mount_key"), connection.Service.Secret)
		},
	})

	return connection
}
