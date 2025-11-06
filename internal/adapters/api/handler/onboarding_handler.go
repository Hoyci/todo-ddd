package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hoyci/todo-ddd/pkg/usecase"
	usecasesetup "github.com/hoyci/todo-ddd/pkg/usecase/setup"
)

type OnboardingHandler struct {
	SetupUC  *usecasesetup.SetupOnboardingUseCase
	Validate *validator.Validate
}

//
// ------------------- Setup -------------------
//

// @Summary Complete user onboarding
// @Description Initiates or finalizes the onboarding process for a new user.
// @Tags Onboarding
// @Accept json
// @Produce json
// @Param onboarding body OnboardingRequest true "Onboarding request payload"
// @Success 200 {object} OnboardingResponse "Onboarding completed successfully"
// @Failure 400 {object} OnboardingErrorResponse "Invalid request payload or validation error"
// @Failure 404 {object} OnboardingErrorResponse "User not found"
// @Failure 409 {object} OnboardingErrorResponse "User with this email already exists"
// @Failure 500 {object} OnboardingErrorResponse "Unexpected internal server error"
// @Router /api/v1/onboarding [post]
func (h *OnboardingHandler) Setup(c *gin.Context) {
	var req OnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.SetupUC.Execute(usecasesetup.SetupOnboardingInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, OnboardingErrorResponse{
				Error: "user with this email already exists",
			})
			return
		case errors.Is(err, usecase.ErrUserNotFound):
			c.JSON(http.StatusNotFound, OnboardingErrorResponse{
				Error: "user not found",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, OnboardingErrorResponse{
				Error: "unexpected error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, OnboardingResponse{
		Message: "onboarding done successfully",
	})
}

//
// ------------------- REQUESTS / RESPONSES -------------------
//

type OnboardingRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

type OnboardingResponse struct {
	Message string `json:"message"`
}

type OnboardingErrorResponse struct {
	Error string `json:"error"`
}
