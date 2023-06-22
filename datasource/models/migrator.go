package models

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&MatrixProfile{},
		&MatrixTransaction{},
		&MatrixApp{},
		&MatrixLibraryItem{},
		&MatrixPost{},
		&MatrixRelease{},
	); err != nil {
		log.Fatal().Err(err).Msg("Error when migrating database")
	}
}
