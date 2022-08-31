package repo

import (
	"entry_task/internal/component"
)

type Repo struct {
	db *component.DB
}

func InitDependencies() *Repo {
	return &Repo{
		db: component.InitDatabase(),
	}
}
