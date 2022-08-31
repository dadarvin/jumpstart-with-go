package handler

import "entry_task/cmd"

type Handler struct {
	user cmd.UserUseCase
}

func New(userUC cmd.UserUseCase) *Handler {
	return &Handler{
		user: userUC,
	}
}
