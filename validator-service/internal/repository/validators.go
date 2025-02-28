package repository

import (
	"gorm.io/gorm"
	"validator-service/internal/models"
)

func CreateValidatorRequest(db *gorm.DB, validator *models.ValidatorRequest) error {
	return db.Create(validator).Error
}

func UpdateValidatorRequest(db *gorm.DB, validator *models.ValidatorRequest) error {
	return db.Save(validator).Error
}

func GetValidatorRequestByUUID(db *gorm.DB, uuid string) (*models.ValidatorRequest, error) {
	var validatorRequest models.ValidatorRequest
	err := db.
		Preload("Keys").
		Where("request_uuid = ?", uuid).
		First(&validatorRequest).
		Error

	return &validatorRequest, err
}

func CreateValidatorKey(db *gorm.DB, validatorKey *models.ValidatorKey) error {
	return db.Create(validatorKey).Error
}
