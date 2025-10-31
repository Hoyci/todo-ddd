package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	_ "github.com/hoyci/todo-ddd/docs/swagger"
	"github.com/hoyci/todo-ddd/internal/adapters/api"
	"github.com/hoyci/todo-ddd/internal/adapters/api/handler"
	"github.com/hoyci/todo-ddd/internal/adapters/db/sqlite"
	usecase "github.com/hoyci/todo-ddd/pkg/usecase/task"
)

func main() {
	db, err := sqlite.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	repo := sqlite.NewSQLiteTaskRepository(db)

	listUC := &usecase.ListTaskUseCase{TaskRepo: repo}
	createUC := &usecase.CreateTaskUseCase{TaskRepo: repo}
	updateUC := &usecase.UpdateTaskUseCase{TaskRepo: repo}
	updateStatusUC := &usecase.UpdateTaskStatusUseCase{TaskRepo: repo}
	deleteUC := &usecase.DeleteTaskUseCase{TaskRepo: repo}

	validate := validator.New()

	taskHandler := &handler.TaskHandler{
		ListUC:         listUC,
		CreateUC:       createUC,
		UpdateUC:       updateUC,
		UpdateStatusUC: updateStatusUC,
		DeleteUC:       deleteUC,
		Validate:       validate,
	}

	router := api.SetupRouter(taskHandler)
	log.Println("Server running on :8080")
	router.Run(":8080")
}
