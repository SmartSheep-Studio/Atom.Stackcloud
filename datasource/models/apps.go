package models

import "gorm.io/datatypes"

type MatrixApp struct {
	Model

	Slug        string                      `json:"slug"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Url         string                      `json:"url"`
	Details     string                      `json:"details"`
	Tags        datatypes.JSONSlice[string] `json:"tags"`
	IsPublished bool                        `json:"is_published"`
	Posts       []MatrixPost                `json:"posts" gorm:"foreignKey:AppID"`
	Releases    []MatrixRelease             `json:"releases" gorm:"foreignKey:AppID"`
	Duplicates  []MatrixLibraryItem         `json:"duplicates" gorm:"foreignKey:AppID"`
	AccountID   uint                        `json:"account_id"`
}
