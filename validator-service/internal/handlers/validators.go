package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
	"validator-service/internal/models"
	"validator-service/internal/repository"
	"validator-service/internal/services"
	"validator-service/internal/utils"
)

const MinNumberOfValidators = 0

const (
	ErrInvalidRequestBody        = "Invalid request body"
	ErrInvalidNumberOfValidators = "Invalid number of validators"
	ErrInvalidFeeRecipient       = "Invalid fee recipient address"
	ErrInternalServer            = "Internal server error"
	ErrCreatingValidator         = "Failed to create validator"
	ErrRequestNotFound           = "Request not found"
	ErrProcessingRequest         = "Error processing request"

	ValidatorCreationInProgress = "Validator creation in progress"
)

type Handler struct {
	db          *gorm.DB
	requestLock sync.Mutex
}

func CreateNewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
	}
}

type CreateValidatorRequest struct {
	NumValidators uint   `json:"num_validators"`
	FeeRecipient  string `json:"fee_recipient"`
}

type CreateValidatorResponse struct {
	RequestId string `json:"request_id"`
	Message   string `json:"message"`
}

type ValidatorStatusResponse struct {
	Status models.RequestStatus `json:"status"`
	Keys   []string             `json:"keys"`
}

type ErrValidatorStatusResponse struct {
	Status  models.RequestStatus `json:"status"`
	Message string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) CreateValidator(c *gin.Context) {
	var req CreateValidatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(ErrInvalidRequestBody, err)
		c.JSON(http.StatusBadRequest, &ErrorResponse{Error: ErrInvalidRequestBody})
		return
	}

	if req.NumValidators <= MinNumberOfValidators {
		log.Println(ErrInvalidNumberOfValidators)
		c.JSON(http.StatusBadRequest, &ErrorResponse{Error: ErrInvalidNumberOfValidators})
		return
	}

	if utils.ValidateAddress(req.FeeRecipient) != true {
		log.Println(ErrInvalidFeeRecipient)
		c.JSON(http.StatusBadRequest, &ErrorResponse{Error: ErrInvalidFeeRecipient})
		return
	}

	validatorRequest := models.ValidatorRequest{
		RequestUUID:   uuid.New().String(),
		NumValidators: req.NumValidators,
		FeeRecipient:  req.FeeRecipient,
		Status:        models.RequestStarted,
	}

	err := repository.CreateValidatorRequest(h.db, &validatorRequest)

	if err != nil {
		log.Printf("%s. %s. Error: %v", ErrCreatingValidator, ErrInternalServer, err)
		c.JSON(http.StatusInternalServerError, &ErrorResponse{Error: ErrInternalServer})
		return
	}

	go services.ProcessValidatorRequest(h.db, &validatorRequest, &h.requestLock)

	c.JSON(http.StatusOK, &CreateValidatorResponse{
		RequestId: validatorRequest.RequestUUID,
		Message:   ValidatorCreationInProgress,
	})
}

func (h *Handler) CheckRequestStatus(c *gin.Context) {
	reqID := c.Param("request_id")
	fmt.Println(reqID)
	validatorRequest, err := repository.GetValidatorRequestByUUID(h.db, reqID)

	if err != nil {
		log.Printf("%s, request_id: %s", ErrRequestNotFound, reqID)
		c.JSON(http.StatusNotFound, &ErrorResponse{Error: ErrRequestNotFound})
		return
	}

	if validatorRequest.Status == models.RequestFailed {
		c.JSON(http.StatusInternalServerError, &ErrValidatorStatusResponse{
			Status:  validatorRequest.Status,
			Message: ErrProcessingRequest,
		})
		return
	}

	c.JSON(http.StatusOK, h.toValidatorStatusResponse(validatorRequest))
}

func (h *Handler) toValidatorStatusResponse(validatorRequest *models.ValidatorRequest) *ValidatorStatusResponse {
	var keys []string

	for _, key := range validatorRequest.Keys {
		keys = append(keys, key.Key)
	}

	return &ValidatorStatusResponse{
		Status: validatorRequest.Status,
		Keys:   keys,
	}
}
