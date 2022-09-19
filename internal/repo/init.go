package repo

import (
	"entry_task/internal/component"
)

type Repo struct {
	db    *component.DB
	cache *component.Cache
}

func InitDependencies() *Repo {
	return &Repo{
		db:    component.InitDatabase(),
		cache: component.InitRedis(),
	}
}
