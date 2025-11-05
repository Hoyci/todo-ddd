package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	usecase "github.com/hoyci/todo-ddd/pkg/usecase/user"
)

// UserHandler agrupa os casos de uso relacionados a usuários
type UserHandler struct {
	CreateUC *usecase.CreateUserUseCase
	UpdateUC *usecase.UpdateUserUseCase
	DeleteUC *usecase.DeleteUserUseCase
	FindUC   *usecase.FindUserUseCase
	Validate *validator.Validate
}

//
// ------------------- CREATE -------------------
//

// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.CreateUC.Execute(usecase.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u := out.User
	c.JSON(http.StatusCreated, UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
}

//
// ------------------- FIND BY ID -------------------
//

// @Summary Get user by ID
// @Description Get a single user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	u, err := h.FindUC.Execute(usecase.FindUserInput{ID: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        u.User.ID,
		Name:      u.User.Name,
		Email:     u.User.Email,
		CreatedAt: u.User.CreatedAt,
		UpdatedAt: u.User.UpdatedAt,
	})
}

//
// ------------------- UPDATE -------------------
//

// @Summary Update a user
// @Description Update name or email of a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UpdateUserRequest true "Updated data"
// @Success 200 {object} UserResponse
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.UpdateUC.Execute(usecase.UpdateUserInput{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u := out.User
	c.JSON(http.StatusOK, UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
}

//
// ------------------- DELETE -------------------
//

// @Summary Delete a user
// @Description Soft delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.DeleteUC.Execute(usecase.DeleteUserInput{ID: id}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

//
// ------------------- REQUESTS / RESPONSES -------------------
//

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
