package models

type MatrixRelease struct {
	Model

	Name         string     `json:"name"`
	Code         string     `json:"code"`
	Description  string     `json:"description"`
	IsPrerelease bool       `json:"is_prerelease"`
	IsPublished  bool       `json:"is_published"`
	Post         MatrixPost `json:"post" gorm:"foreignKey:ReleaseID"`
	AppID        uint       `json:"app_id"`
}
