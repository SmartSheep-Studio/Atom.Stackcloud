package datasource

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDB() *gorm.DB {
	dialector := postgres.Open(viper.GetString("datasource.dsn"))
	conn, err := gorm.Open(dialector, &gorm.Config{Logger: logger.New(&log.Logger, logger.Config{
		Colorful:                  true,
		IgnoreRecordNotFoundError: true,
		LogLevel:                  logger.Info,
	})})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to establish connection with database")
		return nil
	} else {
		return conn
	}
}
