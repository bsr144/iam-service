package mailer

import (
	"context"
	"log"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendOTP(ctx context.Context, email, otp string, expiryMinutes int) error {
	log.Printf(`
========================================
📧 OTP EMAIL (Development Mode)
========================================
To: %s
Subject: Your Verification Code

Your OTP verification code is: %s

This code will expire in %d minutes.

Do not share this code with anyone.
========================================
`, email, otp, expiryMinutes)
	return nil
}

func (s *EmailService) SendWelcome(ctx context.Context, email, firstName string) error {
	log.Printf(`
========================================
📧 WELCOME EMAIL (Development Mode)
========================================
To: %s
Subject: Welcome to Dana Pensiun

Dear %s,

Welcome! Your account has been approved.
You can now log in to access your account.

========================================
`, email, firstName)
	return nil
}

func (s *EmailService) SendPasswordReset(ctx context.Context, email, token string, expiryMinutes int) error {
	log.Printf(`
========================================
📧 PASSWORD RESET EMAIL (Development Mode)
========================================
To: %s
Subject: Password Reset Request

Your password reset token is: %s

This token will expire in %d minutes.

If you did not request this, please ignore this email.
========================================
`, email, token, expiryMinutes)
	return nil
}

func (s *EmailService) SendPINReset(ctx context.Context, email, otp string, expiryMinutes int) error {
	log.Printf(`
========================================
📧 PIN RESET EMAIL (Development Mode)
========================================
To: %s
Subject: PIN Reset Verification Code

Your PIN reset verification code is: %s

This code will expire in %d minutes.

If you did not request this, please secure your account immediately.
========================================
`, email, otp, expiryMinutes)
	return nil
}

func (s *EmailService) SendAdminInvitation(ctx context.Context, email, token string, expiryMinutes int) error {
	log.Printf(`
========================================
📧 ADMIN INVITATION EMAIL (Development Mode)
========================================
To: %s
Subject: You've Been Invited to Join

You have been invited to join as an administrator.

Your invitation token is: %s

This invitation will expire in %d minutes.
========================================
`, email, token, expiryMinutes)
	return nil
}
