package redis

import (
	"context"
	"fmt"
	"log"

	"git.ice.global/packages/beeorm/v4"

	"todo-service/internal/domain/entity"
	beeORMentity "todo-service/internal/repository/beeorm/entity"
	// "todo-service/internal/repository/beeorm/mapper"
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

	fmt.Println("=======start search@@@")
	SearchMyData(s.engine, query)
	fmt.Println("=======end search@@")

	// search := s.engine.GetRedisSearch("todo_cache")
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	page := (offset / limit) + 1
	var entities []*beeORMentity.TodoEntity
	pager := beeorm.NewPager(1, 10)
	searchQuery := beeorm.NewRedisSearchQuery()

	searchQuery.Query("@Description:(" + query + "*)") // Description ফিল্ডে সার্চ করা হচ্ছে

	// if query == "" || query == "*" {
	// 	searchQuery.Query("*")
	// } else {
	// 	searchQuery.FilterString("Description", query)
	// }

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from BeeORM panic: %v\n", r)
		}
	}()

	if s.engine == nil {
		return nil, fmt.Errorf("orm engine is nil")
	}

	fmt.Printf("Debug: Offset=%d, Limit=%d, Query=%s, Page=%+v\n", offset, limit, query, page)

	totalRows := s.engine.RedisSearch(&entities, searchQuery, pager, "")
	fmt.Println("Success! Redis Search found totalRows: ", totalRows)

	var results []*entity.TodoItem
	for _, row := range entities {
		results = append(results, &entity.TodoItem{
			ID:          int(row.ID),
			Description: row.Description,
			// DueDate:     row.DueDate.Format("2006-01-02"),
			// CreatedAt:   row.CreatedAt.Format(time.RFC3339),
		})
	}

	return results, nil
}

func SearchMyData1(engine *beeorm.Engine, searchText string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from SearchMyData panic: %v\n", r)
		}
	}()

	schema := engine.GetRegistry().GetTableSchemaForEntity(&beeORMentity.TodoEntity{})

	search, has := schema.GetRedisSearch(engine)
	if !has {
		log.Fatal("RedisSearch is not configured for this entity")
		return
	}

	indexName := schema.GetTableName()

	fmt.Println("=================SearchMyData indexName:: ", indexName)

	query := &beeorm.RedisSearchQuery{}

	if searchText == "" {
		query.Query("*")
	} else {
		query.Query("@Description:" + searchText + "*")
	}

	// indexName = "e3cf3"
	indexName = "entity.TodoEntity"

	total, keys := search.SearchKeys(indexName, query, beeorm.NewPager(1, 10))

	fmt.Printf("🔍 Total Results: %d\n", total)
	fmt.Printf("🔑 Redis Keys: %v\n", keys)
}

func SearchMyData(engine *beeorm.Engine, searchText string) {
	schema := engine.GetRegistry().GetTableSchemaForEntity(&beeORMentity.TodoEntity{})

	indexName := "entity.TodoEntity"

	search, _ := schema.GetRedisSearch(engine)

	query := &beeorm.RedisSearchQuery{}

	if searchText == "" {
		query.Query("*")
	} else {
		// format: @Description:*sohel* query.Query("@Description:*" + searchText + "*")
		query.Query("@Description:" + searchText + "*")
	}

	pager := beeorm.NewPager(1, 10)

	// var rows []*beeORMentity.TodoEntity

	total, keys := search.Search(indexName, query, pager)

	fmt.Printf("======= start search results =======\n")
	fmt.Printf("🔍 Total Results Found: %d\n", total)
	fmt.Printf("🔑 Redis Keys: %v\n", keys)

	for _, key := range keys {
		fmt.Printf("📌 Found Todo ID Key: %s\n", key)
	}

	// for _, item := range keys {
	// 	fmt.Printf("ID: %d, Description: %s\n", item.ID, item.Description)
	// }
	fmt.Printf("======= end search results =======\n")
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
