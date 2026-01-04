package mapper

import (
	"todo-service/internal/domain/entity"
	beeOrmEntity "todo-service/internal/repository/beeorm/entity"
)

// ToEntity converts BeeORM model to domain entity
func ToEntity(m *beeOrmEntity.TodoEntity) *entity.TodoItem {
	if m == nil {
		return nil
	}

	return &entity.TodoItem{
		ID:          int(m.ID),
		Description: m.Description,
		DueDate:     m.DueDate,
		FileID:      m.FileID,
		CreatedAt:   m.CreatedAt,
	}
}

// ToModel converts domain entity to BeeORM model
func ToModel(todo *entity.TodoItem) *beeOrmEntity.TodoEntity {
	if todo == nil {
		return nil
	}

	return &beeOrmEntity.TodoEntity{
		ID:          uint64(todo.ID),
		Description: todo.Description,
		DueDate:     todo.DueDate,
		FileID:      todo.FileID,
		CreatedAt:   todo.CreatedAt,
	}
}

// ToEntities converts slice of BeeORM models to domain entities
func ToEntities(models []*beeOrmEntity.TodoEntity) []*entity.TodoItem {
	entities := make([]*entity.TodoItem, len(models))
	for i, m := range models {
		entities[i] = ToEntity(m)
	}
	return entities
}

// ToModels converts slice of domain entities to BeeORM models
func ToModels(todos []*entity.TodoItem) []*beeOrmEntity.TodoEntity {
	models := make([]*beeOrmEntity.TodoEntity, len(todos))
	for i, todo := range todos {
		models[i] = ToModel(todo)
	}
	return models
}
