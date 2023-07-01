package models

import "gorm.io/datatypes"

type MatrixRelease struct {
	Model

	Name         string                                   `json:"name"`
	Slug         string                                   `json:"slug"`
	Description  string                                   `json:"description"`
	IsPrerelease bool                                     `json:"is_prerelease"`
	IsPublished  bool                                     `json:"is_published"`
	Post         MatrixPost                               `json:"post" gorm:"foreignKey:ReleaseID"`
	Options      datatypes.JSONType[MatrixReleaseOptions] `json:"options"`
	AppID        uint                                     `json:"app_id"`
}

type MatrixReleaseOptions struct {
	Assets []struct {
		Name         string `json:"name" validate:"required"`
		URL          string `json:"url" validate:"required"`
		Decompressor string `json:"decompressor"`
		Platform     string `json:"platform" validate:"required"`
	} `json:"assets"`
	Preprocessing []struct {
		Name     string `json:"name" validate:"required"`
		Script   string `json:"script" validate:"required"`
		Platform string `json:"platform" validate:"required"`
	} `json:"preprocessing"`
	RunOptions []struct {
		Name     string `json:"name" validate:"required"`
		Script   string `json:"script" validate:"required"`
		Platform string `json:"platform" validate:"required"`
	} `json:"run_options"`
}
