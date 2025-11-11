package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	_ "github.com/hoyci/todo-ddd/docs/swagger"
	"github.com/hoyci/todo-ddd/internal/adapters/api"
	"github.com/hoyci/todo-ddd/internal/adapters/api/handler"
	"github.com/hoyci/todo-ddd/internal/adapters/db/sqlite"
	usecasesetup "github.com/hoyci/todo-ddd/pkg/usecase/setup"
	usecasetask "github.com/hoyci/todo-ddd/pkg/usecase/task"
	usecaseuser "github.com/hoyci/todo-ddd/pkg/usecase/user"
)

func main() {
	db, err := sqlite.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	taskRepo := sqlite.NewSQLiteTaskRepository(db)
	userRepo := sqlite.NewSQLiteUserRepository(db)
	unitOfWork := sqlite.NewSQLiteUnitOfWork(db)

	createTaskUC := &usecasetask.CreateTaskUseCase{UoW: unitOfWork}
	listUC := &usecasetask.ListTaskUseCase{TaskRepo: taskRepo}
	updateUC := &usecasetask.UpdateTaskUseCase{TaskRepo: taskRepo}
	updateStatusUC := &usecasetask.UpdateTaskStatusUseCase{TaskRepo: taskRepo}
	deleteUC := &usecasetask.DeleteTaskUseCase{TaskRepo: taskRepo}

	createUserUC := &usecaseuser.CreateUserUseCase{UserRepo: userRepo}
	updateUserUC := &usecaseuser.UpdateUserUseCase{UserRepo: userRepo}
	deleteUserUC := &usecaseuser.DeleteUserUseCase{UserRepo: userRepo}
	findUserUC := &usecaseuser.FindUserUseCase{UserRepo: userRepo}

	setupUC := &usecasesetup.SetupOnboardingUseCase{UoW: unitOfWork}

	validate := validator.New()

	taskHandler := &handler.TaskHandler{
		ListUC:         listUC,
		CreateUC:       createTaskUC,
		UpdateUC:       updateUC,
		UpdateStatusUC: updateStatusUC,
		DeleteUC:       deleteUC,
		Validate:       validate,
	}

	userHandler := &handler.UserHandler{
		CreateUC: createUserUC,
		UpdateUC: updateUserUC,
		DeleteUC: deleteUserUC,
		FindUC:   findUserUC,
		Validate: validate,
	}

	setupHandler := &handler.OnboardingHandler{
		SetupUC:  setupUC,
		Validate: validate,
	}

	router := api.SetupRouter(taskHandler, userHandler, setupHandler)
	log.Println("Server running on :8080")
	router.Run(":8080")
}
