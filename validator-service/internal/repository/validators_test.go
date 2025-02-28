package repository_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"validator-service/internal/models"
	"validator-service/internal/repository"
)

var baseValidator = models.ValidatorRequest{
	RequestUUID:   "random_uuid",
	NumValidators: 5,
	Status:        models.RequestStarted,
	FeeRecipient:  "0x123",
	Keys: []models.ValidatorKey{
		{Key: "key1", FeeRecipient: "0x123"},
		{Key: "key2", FeeRecipient: "0x124"},
		{Key: "key3", FeeRecipient: "0x125"},
		{Key: "key4", FeeRecipient: "0x126"},
		{Key: "key5", FeeRecipient: "0x127"},
	},
}

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.ValidatorRequest{}, &models.ValidatorKey{})
	return db
}

func TestCreateValidatorRequest(t *testing.T) {
	db := setupTestDB()
	validator := baseValidator

	err := repository.CreateValidatorRequest(db, &validator)
	assert.NoError(t, err)
	assert.NotZero(t, validator.ID)

	var result models.ValidatorRequest
	db.Preload("Keys").First(&result, "random_uuid = ?", validator.RequestUUID)
	assert.Equal(t, baseValidator.RequestUUID, result.RequestUUID)
	assert.Equal(t, baseValidator.NumValidators, result.NumValidators)
	assert.Equal(t, baseValidator.Status, result.Status)
	assert.Equal(t, baseValidator.FeeRecipient, result.FeeRecipient)
	assert.Len(t, result.Keys, 5)
}

func TestUpdateValidatorRequest(t *testing.T) {
	db := setupTestDB()
	validator := baseValidator
	db.Create(&validator)

	validator.Status = models.RequestSuccessful
	err := repository.UpdateValidatorRequest(db, &validator)
	assert.NoError(t, err)

	var result models.ValidatorRequest
	db.Preload("Keys").First(&result, "id = ?", validator.ID)
	assert.Equal(t, models.RequestSuccessful, result.Status)
	assert.Len(t, result.Keys, 5)
}

func TestGetValidatorRequestByUUID(t *testing.T) {
	db := setupTestDB()
	validator := baseValidator
	db.Create(&validator)

	result, err := repository.GetValidatorRequestByUUID(db, validator.RequestUUID)
	assert.NoError(t, err)
	assert.Equal(t, validator.RequestUUID, result.RequestUUID)
	assert.Equal(t, validator.NumValidators, result.NumValidators)
	assert.Equal(t, validator.Status, result.Status)
	assert.Equal(t, validator.FeeRecipient, result.FeeRecipient)
	assert.Len(t, result.Keys, 5)
}

func TestCreateValidatorKey(t *testing.T) {
	db := setupTestDB()
	validator := baseValidator
	db.Create(&validator)

	validatorKey := models.ValidatorKey{
		ValidatorRequestID: validator.ID,
		Key:                "test-key",
		FeeRecipient:       "0x128",
	}
	err := repository.CreateValidatorKey(db, &validatorKey)
	assert.NoError(t, err)

	var result models.ValidatorKey
	db.First(&result, "id = ?", validatorKey.ID)
	assert.Equal(t, validatorKey.Key, result.Key)
	assert.Equal(t, validatorKey.FeeRecipient, result.FeeRecipient)
}
