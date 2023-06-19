package models

type MatrixProfile struct {
	Model

	Nickname   string              `json:"nickname"`
	Experience int64               `json:"experience"`
	Library    []MatrixLibraryItem `json:"library" gorm:"foreignKey:ProfileID"`
	UserID     uint                `json:"user_id"`
}

type MatrixLibraryItem struct {
	Model

	PlayTime  int64 `json:"play_time"`
	ProfileID uint  `json:"profile_id"`
	AppID     uint  `json:"app_id"`
}
