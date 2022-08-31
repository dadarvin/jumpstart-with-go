package usecase

type userRepo interface {
}

type UseCase struct {
	ur userRepo
}

func InitDependencies(ur userRepo) *UseCase {
	return &UseCase{
		ur: ur,
	}
}
