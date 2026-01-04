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

	return results, total, nil
}
