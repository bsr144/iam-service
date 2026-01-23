package contract

import (
	"context"
	"iam-service/internal/role/roledto"
)

type Usecase interface {
	Create(ctx context.Context, req *roledto.CreateRequest) (*roledto.CreateResponse, error)
}
