package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hoyci/todo-ddd/internal/adapters/api/handler"
	swagFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(taskHandler *handler.TaskHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swagFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/tasks", taskHandler.Create)
		v1.PUT("/tasks/:id", taskHandler.Update)
		v1.PATCH("/tasks/:id/status", taskHandler.UpdateStatus)
		v1.GET("/tasks", taskHandler.List)
		v1.DELETE("/tasks/:id", taskHandler.Delete)
	}

	return r
}
