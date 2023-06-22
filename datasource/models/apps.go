package models

import "gorm.io/datatypes"

type MatrixApp struct {
	Model

	Slug         string                                    `json:"slug"`
	Name         string                                    `json:"name"`
	Description  string                                    `json:"description"`
	Url          string                                    `json:"url"`
	Details      string                                    `json:"details"`
	Tags         datatypes.JSONSlice[string]               `json:"tags"`
	IsPublished  bool                                      `json:"is_published"`
	Posts        []MatrixPost                              `json:"posts" gorm:"foreignKey:AppID"`
	Releases     []MatrixRelease                           `json:"releases" gorm:"foreignKey:AppID"`
	Duplicates   []MatrixLibraryItem                       `json:"duplicates" gorm:"foreignKey:AppID"`
	Transactions []MatrixTransaction                       `json:"transactions" gorm:"foreignKey:AppID"`
	PriceOptions datatypes.JSONType[MatrixAppPriceOptions] `json:"price_options"`
	ProfileID    uint                                      `json:"profile_id"`
}

type MatrixAppPriceOptions struct {
	Shop      string `json:"shop" validate:"required"`
	ProductID uint   `json:"product_id" validate:"required"`
	ApiToken  string `json:"api_token" validate:"required"`
}
