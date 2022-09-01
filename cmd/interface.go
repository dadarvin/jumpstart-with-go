package cmd

import "entry_task/internal/model/user"

type UserUseCase interface {
	RegisterUser(user user.User) error
	AuthenticateUser(username string, password string) (interface{}, error)
	UpdateUser(user user.User) error
	GetUserByID(userID int) (user.User, error)
	UploadUserPic(id int, username string, picData string) error
	GetUserPicByID(userID int) (string, error)
}
