package usecase

import "entry_task/internal/model/user"

type userRepo interface {
	UpsertUser(user user.User) error
	GetUser(userID int) (user.User, error)
	UpdateUserPic(picName string, userID int) error
}

type UseCase struct {
	ur userRepo
}

func InitDependencies(ur userRepo) *UseCase {
	return &UseCase{
		ur: ur,
	}
}
