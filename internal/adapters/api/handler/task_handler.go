package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hoyci/todo-ddd/pkg/domain/valueobject"
	usecase "github.com/hoyci/todo-ddd/pkg/usecase/task"
)

type TaskHandler struct {
	CreateUC       *usecase.CreateTaskUseCase
	UpdateUC       *usecase.UpdateTaskUseCase
	UpdateStatusUC *usecase.UpdateTaskStatusUseCase
	DeleteUC       *usecase.DeleteTaskUseCase
	ListUC         *usecase.ListTaskUseCase
	Validate       *validator.Validate
}

//
// ------------------- CREATE -------------------
//

// @Summary Create a new task
// @Description Create a new task for a user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body CreateTaskRequest true "Task data"
// @Success 201 {object} TaskResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.CreateUC.Execute(usecase.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Priority:    valueobject.Priority(req.Priority),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	task := out.Task

	c.JSON(http.StatusCreated, TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    int(task.Priority),
		Status:      string(task.Status),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

//
// ------------------- UPDATE TASK -------------------
//

// @Summary Update a task
// @Description Update title, description or priority
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body UpdateTaskRequest true "Updated data"
// @Success 200 {object} TaskResponse
// @Router /api/v1/tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.UpdateUC.Execute(usecase.UpdateTaskInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Priority:    valueobject.Priority(req.Priority),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    int(task.Priority),
		Status:      string(task.Status),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

//
// ------------------- UPDATE STATUS -------------------
//

// @Summary Update task status
// @Description Update the status of a specific task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param body body UpdateTaskStatusRequest true "Status data"
// @Success 200 {object} TaskResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/tasks/{id}/status [patch]
func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := valueobject.Status(req.Status)
	task, err := h.UpdateStatusUC.Execute(id, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    int(task.Priority),
		Status:      string(task.Status),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

//
// ------------------- LIST -------------------
//

// @Summary List tasks by user
// @Description List all tasks for a given user ID
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {array} TaskResponse
// @Router /api/v1/tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	tasks, err := h.ListUC.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]TaskResponse, 0)
	for _, t := range tasks {
		resp = append(resp, TaskResponse{
			ID:          t.Task.ID,
			Title:       t.Task.Title,
			Description: t.Task.Description,
			Status:      string(t.Task.Status),
			Priority:    int(t.Task.Priority),
			CreatedAt:   t.Task.CreatedAt,
			UpdatedAt:   t.Task.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

//
// ------------------- DELETE -------------------
//

// @Summary Delete a task
// @Description Soft delete a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 204 "No Content"
// @Router /api/v1/tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	_, err := h.DeleteUC.Execute(usecase.DeleteTaskInput{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

//
// ------------------- REQUESTS / RESPONSES -------------------
//

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description"`
	Priority    int    `json:"priority" validate:"required,min=1,max=3"`
}

type TaskResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Priority    int    `json:"priority" validate:"min=1,max=3"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=new in_progress completed"`
}
