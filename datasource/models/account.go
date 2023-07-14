package models

type Account struct {
	Model

	Apps   []App `json:"apps"`
	UserID uint  `json:"user_id"`
}
