package models

import (
	"time"

	"gorm.io/datatypes"
)

// Refer from Quarkpay Transaction Model
type QuarkTransaction struct {
	Model

	Description   string         `json:"description"`
	Amount        float64        `json:"amount"`
	Unit          string         `json:"unit"`
	PayerID       *uint          `json:"payer_id"`
	PayeeID       *uint          `json:"payee_id"`
	Items         datatypes.JSON `json:"items"`
	ExternalLinks datatypes.JSON `json:"external_links"`
	ExpiredAt     *time.Time     `json:"expired_at"`
	ShopID        *uint          `json:"shop_id"`
	IsFinished    bool           `json:"is_finished"`
}

type MatrixTransaction struct {
	Model

	AppID      uint `json:"app_id"`
	QuarkpayID uint `json:"quarkpay_id"`
	ProfileID  uint `json:"profile_id"`
}
