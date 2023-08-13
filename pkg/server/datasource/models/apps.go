package models

import "gorm.io/datatypes"

type App struct {
	Model

	Slug        string                      `json:"slug"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Collections []RecordCollection          `json:"record_collections"`
	Functions   []CloudFunction             `json:"cloud_functions"`
	Tags        datatypes.JSONSlice[string] `json:"tags"`
	AccountID   uint                        `json:"account_id"`
}
