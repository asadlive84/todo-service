package service

import (
	"todo-service/internal/infrastructure/search"

	hitrixService "git.ice.global/packages/hitrix/service"
)

// CustomDIContainer extends hitrix DIContainer
type CustomDIContainer struct {
	*hitrixService.DIContainer
}

// RedisSearch returns Redis Search service
func (d *CustomDIContainer) RedisSearch() *search.RedisSearchService {
	return hitrixService.GetServiceRequired("redis_search").(*search.RedisSearchService)
}

// DI returns custom DI container
func DI() *CustomDIContainer {
	return &CustomDIContainer{
		DIContainer: hitrixService.DI(),
	}
}