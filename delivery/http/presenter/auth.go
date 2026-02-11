package presenter

import (
	"iam-service/delivery/http/dto/response"
	"iam-service/iam/auth/authdto"
)

func ToLoginResponse(resp *authdto.LoginResponse) *response.LoginResponse {
	if resp == nil {
		return nil
	}
	return &response.LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		TokenType:    resp.TokenType,
		User: response.LoginUserResponse{
			ID:       resp.User.ID,
			Email:    resp.User.Email,
			FullName: resp.User.FullName,
			Roles:    resp.User.Roles,
		},
	}
}

func ToRegisterResponse(resp *authdto.RegisterResponse) *response.RegisterResponse {
	if resp == nil {
		return nil
	}
	return &response.RegisterResponse{
		UserID:       resp.UserID,
		Email:        resp.Email,
		Status:       resp.Status,
		OTPExpiresAt: resp.OTPExpiresAt,
	}
}

func ToRegisterSpecialAccountResponse(resp *authdto.RegisterSpecialAccountResponse) *response.RegisterSpecialAccountResponse {
	if resp == nil {
		return nil
	}
	return &response.RegisterSpecialAccountResponse{
		UserID: resp.UserID,
		Email:  resp.Email,
	}
}

func ToVerifyOTPResponse(resp *authdto.VerifyOTPResponse) *response.VerifyOTPResponse {
	if resp == nil {
		return nil
	}
	return &response.VerifyOTPResponse{
		RegistrationToken: resp.RegistrationToken,
		ExpiresIn:         resp.ExpiresIn,
		NextStep:          resp.NextStep,
	}
}

func ToCompleteProfileResponse(resp *authdto.CompleteProfileResponse) *response.CompleteProfileResponse {
	if resp == nil {
		return nil
	}
	return &response.CompleteProfileResponse{
		UserID:   resp.UserID,
		Status:   resp.Status,
		Email:    resp.Email,
		FullName: resp.FullName,
		Message:  resp.Message,
	}
}

func ToResendOTPResponse(resp *authdto.ResendOTPResponse) *response.ResendOTPResponse {
	if resp == nil {
		return nil
	}
	return &response.ResendOTPResponse{
		OTPExpiresAt: resp.OTPExpiresAt,
	}
}

func ToSetupPINResponse(resp *authdto.SetupPINResponse) *response.SetupPINResponse {
	if resp == nil {
		return nil
	}
	return &response.SetupPINResponse{
		PINSetAt: resp.PINSetAt,
	}
}

func ToRequestPasswordResetResponse(resp *authdto.RequestPasswordResetResponse) *response.RequestPasswordResetResponse {
	if resp == nil {
		return nil
	}
	return &response.RequestPasswordResetResponse{
		OTPExpiresAt: resp.OTPExpiresAt,
		EmailMasked:  resp.EmailMasked,
	}
}

func ToResetPasswordResponse(resp *authdto.ResetPasswordResponse) *response.ResetPasswordResponse {
	if resp == nil {
		return nil
	}
	return &response.ResetPasswordResponse{
		Message: resp.Message,
	}
}

func ToInitiateRegistrationResponse(resp *authdto.InitiateRegistrationResponse) *response.InitiateRegistrationResponse {
	if resp == nil {
		return nil
	}
	return &response.InitiateRegistrationResponse{
		RegistrationID: resp.RegistrationID,
		Email:          resp.Email,
		Status:         resp.Status,
		Message:        resp.Message,
		ExpiresAt:      resp.ExpiresAt,
		OTPConfig: response.InitiateRegistrationOTPConfig{
			ExpiresInMinutes:      resp.OTPConfig.ExpiresInMinutes,
			ResendCooldownSeconds: resp.OTPConfig.ResendCooldownSeconds,
		},
	}
}

func ToVerifyRegistrationOTPResponse(resp *authdto.VerifyRegistrationOTPResponse) *response.VerifyRegistrationOTPResponse {
	if resp == nil {
		return nil
	}
	return &response.VerifyRegistrationOTPResponse{
		RegistrationID:    resp.RegistrationID,
		Status:            resp.Status,
		Message:           resp.Message,
		RegistrationToken: resp.RegistrationToken,
		TokenExpiresAt:    resp.TokenExpiresAt,
		NextStep: response.VerifyRegistrationOTPNextStep{
			Action:   resp.NextStep.Action,
			Endpoint: resp.NextStep.Endpoint,
		},
	}
}

func ToResendRegistrationOTPResponse(resp *authdto.ResendRegistrationOTPResponse) *response.ResendRegistrationOTPResponse {
	if resp == nil {
		return nil
	}
	return &response.ResendRegistrationOTPResponse{
		RegistrationID:        resp.RegistrationID,
		Message:               resp.Message,
		ExpiresAt:             resp.ExpiresAt,
		ResendsRemaining:      resp.ResendsRemaining,
		NextResendAvailableAt: resp.NextResendAvailableAt,
	}
}

func ToRegistrationStatusResponse(resp *authdto.RegistrationStatusResponse) *response.RegistrationStatusResponse {
	if resp == nil {
		return nil
	}
	return &response.RegistrationStatusResponse{
		RegistrationID:       resp.RegistrationID,
		Email:                resp.Email,
		Status:               resp.Status,
		ExpiresAt:            resp.ExpiresAt,
		OTPAttemptsRemaining: resp.OTPAttemptsRemaining,
		ResendsRemaining:     resp.ResendsRemaining,
	}
}

func ToCompleteRegistrationResponse(resp *authdto.CompleteRegistrationResponse) *response.CompleteRegistrationResponse {
	if resp == nil {
		return nil
	}
	return &response.CompleteRegistrationResponse{
		UserID:  resp.UserID,
		Email:   resp.Email,
		Status:  resp.Status,
		Message: resp.Message,
		Profile: response.CompleteRegistrationProfile{
			FirstName: resp.Profile.FirstName,
			LastName:  resp.Profile.LastName,
		},
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		TokenType:    resp.TokenType,
		ExpiresIn:    resp.ExpiresIn,
	}
}
