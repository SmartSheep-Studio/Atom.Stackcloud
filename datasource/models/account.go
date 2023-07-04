package models

type MatrixAccount struct {
	Model

	Nickname   string              `json:"nickname"`
	Experience int64               `json:"experience"`
	Library    []MatrixLibraryItem `json:"library" gorm:"foreignKey:AccountID"`
	UserID     uint                `json:"user_id"`
}

type MatrixLibraryItem struct {
	Model

	PlayTime  int64           `json:"play_time"`
	CloudSave MatrixCloudSave `json:"cloud_save" gorm:"foreignKey:LibraryID"`
	AccountID uint            `json:"account_id"`
	AppID     uint            `json:"app_id"`
}
