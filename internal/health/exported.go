package health

import "iam-service/internal/health/internal"

type Usecase interface {
	CheckHealth() error
}

func NewUsecase() Usecase {
	return internal.NewUsecase()
}
