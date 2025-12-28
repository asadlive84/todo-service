package search

import (
	"context"
	"fmt"
	"strings"
	"time"
	"todo-service/internal/infrastructure/helper"

	"git.ice.global/packages/beeorm/v4"
	"git.ice.global/packages/hitrix/service"
)

type RedisSearchService struct {
	redisCache *beeorm.RedisCache
}

func NewRedisSearchService() *RedisSearchService {
	ormEngine := service.DI().OrmEngine()
	redisCache := ormEngine.GetRedis()

	return &RedisSearchService{
		redisCache: redisCache,
	}
}

// Create Todo search index
func (s *RedisSearchService) CreateTodoIndex(ctx context.Context) error {
	script := `
		local result = redis.call('FT.CREATE', 'idx:todos',
			'ON', 'HASH',
			'PREFIX', '1', 'todo:',
			'SCHEMA',
			'id', 'NUMERIC', 'SORTABLE',
			'description', 'TEXT', 'WEIGHT', '2.0',
			'fileid', 'TEXT',
			'dueDate', 'NUMERIC', 'SORTABLE',
			'createdAt', 'NUMERIC', 'SORTABLE'
		)
		return result
	`

	// Use defer to recover from panic
	defer func() {
		if r := recover(); r != nil {
			// Check if error is about index already existing
			errMsg := fmt.Sprintf("%v", r)
			if !strings.Contains(errMsg, "Index already exists") {
				// Re-panic if it's a different error
				panic(r)
			}
			// Otherwise, silently ignore "Index already exists" error
		}
	}()

	s.redisCache.Eval(script, []string{})

	return nil
}

func (s *RedisSearchService) IndexTodo(ctx context.Context, todoID uint64, description, fileID string, dueDate, createdAt time.Time) error {
	key := fmt.Sprintf("todo:%d", todoID)

	s.redisCache.HSet(key,
		"id", todoID,
		"description", description,
		"fileid", fileID,
		"dueDate", dueDate.Unix(),
		"createdAt", createdAt.Unix(),
	)

	return nil
}

// Search todos by query
func (s *RedisSearchService) SearchTodos(ctx context.Context, query string, offset, limit int) ([]map[string]interface{}, int64, error) {
	script := fmt.Sprintf(`
		return redis.call('FT.SEARCH', 'idx:todos', '%s', 'LIMIT', '%d', '%d')
	`, query, offset, limit)

	result := s.redisCache.Eval(script, []string{})

	return s.parseSearchResults(result)
}

// Parse search results
func (s *RedisSearchService) parseSearchResults(result interface{}) ([]map[string]interface{}, int64, error) {
	// 1. Assert the top-level response is a slice
	data, ok := result.([]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("response is not a slice")
	}

	// 2. Extract Count (Index 0)
	// If the slice is empty or too short, return empty
	if len(data) == 0 {
		return []map[string]interface{}{}, 0, nil
	}

	countVal, _ := data[0].(int64) 
	if c, ok := data[0].(int); ok {
		countVal = int64(c)
	}

	parsedResults := make([]map[string]interface{}, 0)

	for i := 1; i < len(data); i += 2 {
		if i+1 >= len(data) {
			break 
		}

		fieldsRaw, ok := data[i+1].([]interface{})
		if !ok {
			fmt.Printf("Error: Item at index %d is not a field list\n", i+1)
			continue
		}

		// 4. Convert the Flat List ["id", "8", "desc", "text"] into a Map
		itemMap := make(map[string]interface{})

		for k := 0; k < len(fieldsRaw); k += 2 {
			if k+1 >= len(fieldsRaw) {
				break
			}

			// Redis often returns []byte, so we use a helper to stringify keys/values
			key := helper.ToString(fieldsRaw[k])
			val := helper.ToString(fieldsRaw[k+1])

			itemMap[key] = val
		}

		parsedResults = append(parsedResults, itemMap)
	}

	return parsedResults, countVal, nil
}



// // Search by description
// func (s *RedisSearchService) SearchByDescription(ctx context.Context, description string) ([]map[string]interface{}, error) {
// 	query := fmt.Sprintf("@description:%s", description)
// 	todos, _, err := s.SearchTodos(ctx, query, 0, 10)
// 	return todos, err
// }

// // Delete todo from index
// func (s *RedisSearchService) DeleteTodo(ctx context.Context, todoID uint64) error {
// 	key := fmt.Sprintf("todo:%d", todoID)
// 	s.redisCache.Del(key)
// 	return nil
// }

// // Get index info
// func (s *RedisSearchService) GetIndexInfo(ctx context.Context) (map[string]interface{}, error) {
// 	script := `
// 		return redis.call('FT.INFO', 'idx:todos')
// 	`

// 	result := s.redisCache.Eval(script, []string{})

// 	return s.parseIndexInfo(result)
// }

// // Parse index info
// func (s *RedisSearchService) parseIndexInfo(result interface{}) (map[string]interface{}, error) {
// 	if result == nil {
// 		return nil, fmt.Errorf("no index info found")
// 	}

// 	info := make(map[string]interface{})

// 	if resultSlice, ok := result.([]interface{}); ok {
// 		for i := 0; i < len(resultSlice); i += 2 {
// 			if i+1 < len(resultSlice) {
// 				key, ok1 := resultSlice[i].(string)
// 				value := resultSlice[i+1]

// 				if ok1 {
// 					info[key] = value
// 				}
// 			}
// 		}
// 	}

// 	return info, nil
// }

// // Drop index (for testing)
// func (s *RedisSearchService) DropIndex(ctx context.Context) error {
// 	script := `
// 		return redis.call('FT.DROPINDEX', 'idx:todos', 'DD')
// 	`

// 	s.redisCache.Eval(script, []string{})
// 	return nil
// }
