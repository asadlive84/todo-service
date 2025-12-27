package service

import (
	hitrixService "git.ice.global/packages/hitrix/service"
	"github.com/sarulabs/di"
	"todo-service/internal/infrastructure/search"
)

// ServiceProviderRedisSearch creates Redis Search service
func ServiceProviderRedisSearch() *hitrixService.DefinitionGlobal {
	return &hitrixService.DefinitionGlobal{
		Name: "redis_search",
		Build: func(ctn di.Container) (interface{}, error) {
			return search.NewRedisSearchService(), nil
		},
	}
}