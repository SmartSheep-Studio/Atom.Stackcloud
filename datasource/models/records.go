package models

import "gorm.io/datatypes"

type RecordCollection struct {
	Model

	Slug        string                      `json:"slug" gorm:"uniqueIndex"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Tags        datatypes.JSONSlice[string] `json:"tags"`
	Records     []Record                    `json:"records" gorm:"foreignKey:CollectionID"`
	AppID       uint                        `json:"app_id"`
}

type Record struct {
	Model

	Payload      datatypes.JSON `json:"payload"`
	CollectionID uint           `json:"collection_id"`
}
