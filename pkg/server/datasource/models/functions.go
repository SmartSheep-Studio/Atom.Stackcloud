package models

import "gorm.io/datatypes"

type CloudFunction struct {
	Model

	Slug        string                      `json:"slug" gorm:"uniqueIndex"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Script      string                      `json:"script"`
	Tags        datatypes.JSONSlice[string] `json:"tags"`
	AppID       uint                        `json:"app_id"`
}
