package redis

import (
	"context"
	"fmt"

	"git.ice.global/packages/beeorm/v4"

	"todo-service/internal/domain/entity"
	beeORMentity "todo-service/internal/repository/beeorm/entity"
	"todo-service/internal/repository/beeorm/mapper"
)

type RedisSearchService struct {
	engine *beeorm.Engine
}

func NewRedisSearchService(engine *beeorm.Engine) *RedisSearchService {
	return &RedisSearchService{
		engine: engine,
	}
}

// Search todos by query
func (s *RedisSearchService) SearchTodos(ctx context.Context, query string, offset, limit int) ([]*entity.TodoItem, error) {
	searchQuery := beeorm.NewRedisSearchQuery()

	searchQuery.Query("@Description:(" + query + ")")

	// var todos []*entity.TodoItem

	pager := beeorm.NewPager(offset, limit)

	var models []*beeORMentity.TodoEntity

	totalRows := s.engine.RedisSearch(models, searchQuery, pager)

	fmt.Printf("Total found: %d\n", totalRows)

	if models != nil {
		// Convert to domain entities
		entities := make([]*entity.TodoItem, len(models))

		for i, m := range models {
			entities[i] = mapper.ToEntity(m)
		}
	}

	return nil, fmt.Errorf("slice is nil or no data found %+v", "issue####")

}

// Create Todo search index
// func (s *RedisSearchService) CreateTodoIndex(ctx context.Context) error {
// 	script := `
// 		local result = redis.call('FT.CREATE', 'idx:todos',
// 			'ON', 'HASH',
// 			'PREFIX', '1', 'todo:',
// 			'SCHEMA',
// 			'id', 'NUMERIC', 'SORTABLE',
// 			'description', 'TEXT', 'WEIGHT', '2.0',
// 			'fileid', 'TEXT',
// 			'dueDate', 'NUMERIC', 'SORTABLE',
// 			'createdAt', 'NUMERIC', 'SORTABLE'
// 		)
// 		return result
// 	`

// 	// Use defer to recover from panic
// 	defer func() {
// 		if r := recover(); r != nil {
// 			// Check if error is about index already existing
// 			errMsg := fmt.Sprintf("%v", r)
// 			if !strings.Contains(errMsg, "Index already exists") {
// 				// Re-panic if it's a different error
// 				panic(r)
// 			}
// 			// Otherwise, silently ignore "Index already exists" error
// 		}
// 	}()

// 	s.redisCache.Eval(script, []string{})

// 	return nil
// }

// func (s *RedisSearchService) IndexTodo(ctx context.Context, todoID uint64, description, fileID string, dueDate, createdAt time.Time) error {
// 	key := fmt.Sprintf("todo:%d", todoID)

// 	s.redisCache.HSet(key,
// 		"id", todoID,
// 		"description", description,
// 		"fileid", fileID,
// 		"dueDate", dueDate.Unix(),
// 		"createdAt", createdAt.Unix(),
// 	)

// 	return nil
// }

// Search todos by query
// func (s *RedisSearchService) SearchTodos(ctx context.Context, query string, offset, limit int) ([]map[string]interface{}, int64, error) {
// 	script := fmt.Sprintf(`
// 		return redis.call('FT.SEARCH', 'idx:todos', '%s', 'LIMIT', '%d', '%d')
// 	`, query, offset, limit)

// 	result := s.redisCache.Eval(script, []string{})

// 	return s.parseSearchResults(result)
// }

// Parse search results
// func (s *RedisSearchService) parseSearchResults(result interface{}) ([]map[string]interface{}, int64, error) {
// 	// 1. Assert the top-level response is a slice
// 	data, ok := result.([]interface{})
// 	if !ok {
// 		return nil, 0, fmt.Errorf("response is not a slice")
// 	}

// 	// 2. Extract Count (Index 0)
// 	// If the slice is empty or too short, return empty
// 	if len(data) == 0 {
// 		return []map[string]interface{}{}, 0, nil
// 	}

// 	countVal, _ := data[0].(int64)
// 	if c, ok := data[0].(int); ok {
// 		countVal = int64(c)
// 	}

// 	parsedResults := make([]map[string]interface{}, 0)

// 	for i := 1; i < len(data); i += 2 {
// 		if i+1 >= len(data) {
// 			break
// 		}

// 		fieldsRaw, ok := data[i+1].([]interface{})
// 		if !ok {
// 			fmt.Printf("Error: Item at index %d is not a field list\n", i+1)
// 			continue
// 		}

// 		// 4. Convert the Flat List ["id", "8", "desc", "text"] into a Map
// 		itemMap := make(map[string]interface{})

// 		for k := 0; k < len(fieldsRaw); k += 2 {
// 			if k+1 >= len(fieldsRaw) {
// 				break
// 			}

// 			// Redis often returns []byte, so we use a helper to stringify keys/values
// 			key := helper.ToString(fieldsRaw[k])
// 			val := helper.ToString(fieldsRaw[k+1])

// 			itemMap[key] = val
// 		}

// 		parsedResults = append(parsedResults, itemMap)
// 	}

// 	return parsedResults, countVal, nil
// }
