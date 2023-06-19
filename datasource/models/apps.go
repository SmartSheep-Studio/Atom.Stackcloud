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
	UserID      uint                        `json:"app_id"`
}
