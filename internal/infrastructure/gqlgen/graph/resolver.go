package graph

//go:generate go tool gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

import "todo-service/internal/usecase"

type Resolver struct {
	TodoUseCase *usecase.TodoUseCase
}
