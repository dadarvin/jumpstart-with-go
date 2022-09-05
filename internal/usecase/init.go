package usecase

import (
	"database/sql"
	"entry_task/internal/model/user"
	"image"
)

type userRepo interface {
	CreateTx() (*sql.Tx, error)
	GetJWT(username string) (string, error)
	UpsertUser(user user.User, tx *sql.Tx) error
	GetUserByID(userID int) (user.User, error)
	GetUserByName(username string) (user.User, error)
	UpdateUserPic(picName string, userID int, tx *sql.Tx) error
	UploadUserPic(image image.Image, fileName string) error
}

type UseCase struct {
	ur userRepo
}

func InitDependencies(ur userRepo) *UseCase {
	return &UseCase{
		ur: ur,
	}
}
