package services

import (
	"gorm.io/gorm"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
	"validator-service/internal/models"
	"validator-service/internal/repository"
)

const DelayTime = 20 * time.Millisecond

const (
	ErrCreatingValidator              = "Failed to create validator"
	ErrCreatingValidatorKey           = "Failed to create validator key"
	ErrGeneratingRandomString         = "Failed to generate random string"
	ErrUpdatingValidatorRequestStatus = "Failed to update validator request status"
)

type Result struct {
	key string
	err error
}

func ProcessValidatorRequest(db *gorm.DB, validatorRequest *models.ValidatorRequest, requestLock *sync.Mutex) {
	var keys []string
	var errors []error
	var wg sync.WaitGroup
	var keyLock sync.Mutex

	for i := uint(0); i < validatorRequest.NumValidators; i++ {
		wg.Add(1)
		go createValidator(&keys, &errors, &wg, &keyLock)
	}

	wg.Wait() // wait for all validators to be created

	if len(errors) > 0 {
		log.Printf("%s: %v", ErrCreatingValidator, errors)
		updateValidatorStatus(db, validatorRequest, models.RequestFailed, requestLock)

		return
	}

	for i := 0; i < len(keys); i++ {
		validatorKey := models.ValidatorKey{
			ValidatorRequestID: validatorRequest.ID,
			Key:                keys[i],
			FeeRecipient:       validatorRequest.FeeRecipient,
		}
		err := repository.CreateValidatorKey(db, &validatorKey)
		if err != nil {
			log.Printf("%s: %v", ErrCreatingValidatorKey, err)
			updateValidatorStatus(db, validatorRequest, models.RequestFailed, requestLock)

			return
		}
	}

	updateValidatorStatus(db, validatorRequest, models.RequestSuccessful, requestLock)
}

func createValidator(keys *[]string, errs *[]error, wg *sync.WaitGroup, keyLock *sync.Mutex) {
	defer wg.Done()

	time.Sleep(DelayTime)

	keyLock.Lock()
	key, err := generateRandomString(32)
	keyLock.Unlock()

	if err != nil {
		log.Printf("%s: %v", ErrGeneratingRandomString, err)
		*errs = append(*errs, err)
		return
	}

	*keys = append(*keys, key)
}

func updateValidatorStatus(db *gorm.DB, validatorRequest *models.ValidatorRequest, status models.RequestStatus, lock *sync.Mutex) {
	validatorRequest.Status = status

	lock.Lock()
	err := repository.UpdateValidatorRequest(db, validatorRequest)
	lock.Unlock()

	if err != nil {
		log.Printf("%s, requestID: %s, error: %v", ErrUpdatingValidatorRequestStatus, validatorRequest.RequestUUID, err)
	}
}

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	randNum := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		sb.WriteByte(charset[randNum.Intn(len(charset))])
	}

	return sb.String(), nil
}
