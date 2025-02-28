package models

import "gorm.io/gorm"

type RequestStatus string

const (
	RequestStarted    RequestStatus = "started"
	RequestSuccessful RequestStatus = "successful"
	RequestFailed     RequestStatus = "failed"
)

type ValidatorRequest struct {
	gorm.Model
	RequestUUID   string         `json:"request_uuid"`
	NumValidators uint           `json:"num_validators"`
	FeeRecipient  string         `json:"fee_recipient"`
	Status        RequestStatus  `json:"status"`
	Keys          []ValidatorKey `json:"keys" gorm:"foreignKey:ValidatorRequestID"`
}

type ValidatorKey struct {
	gorm.Model
	ValidatorRequestID uint   `json:"validator_request_id"`
	Key                string `json:"key"`
	FeeRecipient       string `json:"fee_recipient"`
}
