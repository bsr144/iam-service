package health

import "iam-service/iam/health/internal"

type Usecase interface {
	CheckHealth() error
}

func NewUsecase() Usecase {
	return internal.NewUsecase()
}
