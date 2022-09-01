package usecase

import (
	"database/sql"
	"entry_task/internal/model/user"
)

type userRepo interface {
	CreateTx() (*sql.Tx, error)
	UpsertUser(user user.User, tx *sql.Tx) error
	GetUserByID(userID int) (user.User, error)
	GetUserByName(username string) (user.User, error)
	UpdateUserPic(picName string, userID int, tx *sql.Tx) error
}

type UseCase struct {
	ur userRepo
}

func InitDependencies(ur userRepo) *UseCase {
	return &UseCase{
		ur: ur,
	}
}
