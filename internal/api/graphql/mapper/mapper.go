package mapper

import (
	"fmt"
	"time"
	"todo-service/internal/domain/dto"
	"todo-service/internal/domain/entity"
)

// InputToEntity converts CreateTodoInput (DTO) to TodoItem (Domain Entity)
func InputToEntity(input dto.CreateTodoInput) (*entity.TodoItem, error) {
	// Parse DueDate from ISO 8601 string to time.Time
	var dueDate time.Time
	if input.DueDate != "" {
		var err error
		dueDate, err = time.Parse(time.RFC3339, input.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid DueDate format: %w", err)
		}
	}

	return &entity.TodoItem{
		Description: input.Description,
		DueDate:     dueDate.UTC(),
		FileID:      input.FileID,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

// EntityToPayload converts TodoItem (Domain Entity) to CreateTodoPayload (DTO)
func EntityToPayload(todoItem *entity.TodoItem) *dto.CreateTodoPayload {

	fmt.Println("======todoItem======>", todoItem)
	return &dto.CreateTodoPayload{
		Todo: &dto.Todo{
			ID:          int32(todoItem.ID),
			Description: todoItem.Description,
			DueDate:     todoItem.DueDate.Format(time.RFC3339),
			FileID:      todoItem.FileID,
			CreatedAt:   todoItem.CreatedAt.Format(time.RFC3339),
		},
		Error: nil, // or &dto.FieldError{} if you need empty error
	}
}

// EntityToPayloadWithError converts TodoItem with error message
func EntityToPayloadWithError(todoItem *entity.TodoItem, errMsg string, field string) *dto.CreateTodoPayload {
	payload := &dto.CreateTodoPayload{
		Error: &dto.FieldError{
			Field:   field,
			Message: errMsg,
		},
	}

	if todoItem != nil {
		payload.Todo = &dto.Todo{
			ID:          int32(todoItem.ID),
			Description: todoItem.Description,
			DueDate:     todoItem.DueDate.Format(time.RFC3339),
			FileID:      todoItem.FileID,
			CreatedAt:   todoItem.CreatedAt.Format(time.RFC3339),
		}
	}

	return payload
}

func EntitiesToSearchResult(items []*entity.TodoItem, total, offset, limit int) *dto.SearchResult {
	// Convert each TodoItem to Todo DTO
	todos := make([]*dto.Todo, len(items))
	for i, item := range items {
		todos[i] = EntityToTodo(item)
	}

	return &dto.SearchResult{
		Results: todos,
		Total:   int32(total),
		Offset:  int32(offset),
		Limit:   int32(limit),
	}
}

// EntityToTodo converts single TodoItem entity to Todo DTO
func EntityToTodo(item *entity.TodoItem) *dto.Todo {
	return &dto.Todo{
		ID:          int32(item.ID),
		Description: item.Description,
		DueDate:     item.DueDate.Format(time.RFC3339),
		FileID:      item.FileID,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339), // if you have UpdatedAt
	}
}
