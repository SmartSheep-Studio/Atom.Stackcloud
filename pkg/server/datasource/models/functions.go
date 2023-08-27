package models

import "gorm.io/datatypes"

type CloudFunction struct {
	Model

	Slug        string                      `json:"slug" gorm:"index:cloud_functions_pkey,unique"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Script      string                      `json:"script"`
	Tags        datatypes.JSONSlice[string] `json:"tags"`
	AppID       uint                        `json:"app_id" gorm:"index:cloud_functions_pkey,unique"`
}
