package _dev

// Repo mockgen
//go:generate mockgen -source=../internal/repo/init.go -destination=../internal/repo/repo_mock_test.go -package=repo

// Usecase mockgen
//go:generate mockgen -source=../internal/usecase/init.go -destination=../internal/usecase/user_mock_test.go -package=usecase