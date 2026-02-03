package contract

import (
	"context"

	"iam-service/iam/auth/authdto"
)

type RegistrationUseCase interface {
	Register(ctx context.Context, req *authdto.RegisterRequest) (*authdto.RegisterResponse, error)
	VerifyOTP(ctx context.Context, req *authdto.VerifyOTPRequest) (*authdto.VerifyOTPResponse, error)
	CompleteProfile(ctx context.Context, req *authdto.CompleteProfileRequest) (*authdto.CompleteProfileResponse, error)
	ResendOTP(ctx context.Context, req *authdto.ResendOTPRequest) (*authdto.ResendOTPResponse, error)
}
