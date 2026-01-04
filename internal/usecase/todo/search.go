package todo

import (
	"context"
	"fmt"
	domain "todo-service/internal/domain/entity"
)

func (uc *TodoUseCase) Search(ctx context.Context, query string, offset, limit int32) (res []*domain.TodoItem, total int64, err error) {

	off := 0
	lim := 10

	if offset != 0 {
		off = int(offset)
	}
	if limit != 0 {
		lim = int(limit)
	}

	fmt.Println("=================Search Use case=======================")

	results, err := uc.search.SearchTodos(ctx, query, off, lim)
	if err != nil {
		fmt.Printf("error comes from %+v\n", err)
		return nil, 0, err
	}

	total = int64(len(results))

	fmt.Printf("Search results: %+v\n", results)
	fmt.Println("==============Search END==========================")

	// todos := make([]*domain.TodoItem, 0, len(results))
	// for _, result := range results {
	// 	// Safe type assertion for ID
	// 	id := 0
	// 	if idStr, ok := result["id"].(string); ok {
	// 		if parsedID, err := strconv.Atoi(idStr); err == nil {
	// 			id = parsedID
	// 		}
	// 	}

	// 	// Safe type assertion for Description
	// 	description := ""
	// 	if desc, ok := result["description"].(string); ok {
	// 		description = desc
	// 	}

	// 	// Safe type assertion for FileID
	// 	fileID := ""
	// 	if fid, ok := result["fileid"].(string); ok {
	// 		fileID = fid
	// 	}

	// 	// Safe type assertion for DueDate
	// 	var dueDate time.Time
	// 	switch v := result["dueDate"].(type) {
	// 	case time.Time:
	// 		dueDate = v
	// 	case string:
	// 		// Try parsing if it's a string
	// 		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
	// 			dueDate = parsed
	// 		}
	// 	case int64:
	// 		// Unix timestamp
	// 		dueDate = time.Unix(v, 0)
	// 	}

	// 	// Safe type assertion for CreatedAt
	// 	var createdAt time.Time
	// 	switch v := result["createdAt"].(type) {
	// 	case time.Time:
	// 		createdAt = v
	// 	case string:
	// 		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
	// 			createdAt = parsed
	// 		}
	// 	case int64:
	// 		createdAt = time.Unix(v, 0)
	// 	}

	// 	// Safe type assertion for UpdatedAt (if exists)
	// 	var updatedAt time.Time
	// 	if val, exists := result["updatedAt"]; exists {
	// 		switch v := val.(type) {
	// 		case time.Time:
	// 			updatedAt = v
	// 		case string:
	// 			if parsed, err := time.Parse(time.RFC3339, v); err == nil {
	// 				updatedAt = parsed
	// 			}
	// 		case int64:
	// 			updatedAt = time.Unix(v, 0)
	// 		}
	// 	}

	// 	todo := &domain.TodoItem{
	// 		ID:          id,
	// 		Description: description,
	// 		DueDate:     dueDate,
	// 		FileID:      fileID,
	// 		CreatedAt:   createdAt,
	// 		UpdatedAt:   updatedAt,
	// 	}

	// 	todos = append(todos, todo)
	// }

	// fmt.Println("==============HI==========================")

	return results, total, nil
}
