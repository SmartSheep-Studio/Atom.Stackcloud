package models

import "gorm.io/datatypes"

type LineupApp struct {
	Model

	Slug        string                      `json:"slug"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Url         string                      `json:"url"`
	Details     string                      `json:"details"`
	Tags        datatypes.JSONSlice[string] `json:"tags"`
	UserID      uint                        `json:"app_id"`
}
