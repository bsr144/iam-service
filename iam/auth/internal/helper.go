package internal

import (
	"crypto/rand"
	"iam-service/pkg/errors"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) validatePassword(password string) error {
	if len(password) < PasswordMinLength {
		return errors.ErrValidation("Password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case char >= '!' && char <= '/' || char >= ':' && char <= '@' || char >= '[' && char <= '`' || char >= '{' && char <= '~':
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.ErrValidation("Password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.ErrValidation("Password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.ErrValidation("Password must contain at least one number")
	}
	if !hasSpecial {
		return errors.ErrValidation("Password must contain at least one special character")
	}

	return nil
}

func (uc *usecase) generateOTP() (otp string, otpHash string, err error) {
	digits := make([]byte, OTPLength)
	for i := 0; i < OTPLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", "", err
		}
		digits[i] = byte('0' + n.Int64())
	}
	otp = string(digits)

	hash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return otp, string(hash), nil
}

func (uc *usecase) generateRegistrationToken(userID, tenantID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID.String(),
		"tenant_id": tenantID.String(),
		"purpose":   "registration",
		"exp":       time.Now().Add(10 * time.Minute).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.Config.JWT.AccessSecret))
}

func (uc *usecase) parseRegistrationToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrTokenInvalid()
		}
		return []byte(uc.Config.JWT.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrTokenInvalid()
	}

	if purpose, ok := claims["purpose"].(string); !ok || purpose != "registration" {
		return nil, errors.ErrTokenInvalid()
	}

	return claims, nil
}
