package errors

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	CodeInternal           = "ERR_INTERNAL"
	CodeValidation         = "ERR_VALIDATION"
	CodeNotFound           = "ERR_NOT_FOUND"
	CodeConflict           = "ERR_CONFLICT"
	CodeBadRequest         = "ERR_BAD_REQUEST"
	CodeUnauthorized       = "ERR_UNAUTHORIZED"
	CodeForbidden          = "ERR_FORBIDDEN"
	CodeTooManyRequests    = "ERR_TOO_MANY_REQUESTS"
	CodeServiceUnavailable = "ERR_SERVICE_UNAVAILABLE"

	CodeInvalidCredentials = "ERR_INVALID_CREDENTIALS"
	CodeTokenExpired       = "ERR_TOKEN_EXPIRED"
	CodeTokenInvalid       = "ERR_TOKEN_INVALID"
	CodeSessionExpired     = "ERR_SESSION_EXPIRED"
	CodeOTPInvalid         = "ERR_OTP_INVALID"
	CodeOTPExpired         = "ERR_OTP_EXPIRED"
	CodePINInvalid         = "ERR_PIN_INVALID"
	CodePINLocked          = "ERR_PIN_LOCKED"
	CodePINRequired        = "ERR_PIN_REQUIRED"

	CodeUserNotFound      = "ERR_USER_NOT_FOUND"
	CodeUserAlreadyExists = "ERR_USER_ALREADY_EXISTS"
	CodeUserNotApproved   = "ERR_USER_NOT_APPROVED"
	CodeUserSuspended     = "ERR_USER_SUSPENDED"
	CodeUserInactive      = "ERR_USER_INACTIVE"
	CodeEmailNotVerified  = "ERR_EMAIL_NOT_VERIFIED"
	CodeProfileIncomplete = "ERR_PROFILE_INCOMPLETE"

	CodeTenantNotFound  = "ERR_TENANT_NOT_FOUND"
	CodeTenantInactive  = "ERR_TENANT_INACTIVE"
	CodeTenantSuspended = "ERR_TENANT_SUSPENDED"
	CodeInvalidTenant   = "ERR_INVALID_TENANT"

	CodePermissionDenied      = "ERR_PERMISSION_DENIED"
	CodeRoleNotFound          = "ERR_ROLE_NOT_FOUND"
	CodeInvalidPermission     = "ERR_INVALID_PERMISSION"
	CodeAccessForbidden       = "ERR_ACCESS_FORBIDDEN"
	CodePlatformAdminRequired = "ERR_PLATFORM_ADMIN_REQUIRED"

	CodeEmployeeNotFound = "ERR_EMPLOYEE_NOT_FOUND"
	CodeEmployeeExists   = "ERR_EMPLOYEE_EXISTS"
	CodeInvalidNIK       = "ERR_INVALID_NIK"

	CodeContributionNotFound    = "ERR_CONTRIBUTION_NOT_FOUND"
	CodeContributionInvalid     = "ERR_CONTRIBUTION_INVALID"
	CodeContributionDuplicate   = "ERR_CONTRIBUTION_DUPLICATE"
	CodeInvalidContributionData = "ERR_INVALID_CONTRIBUTION_DATA"

	CodeAllocationNotFound = "ERR_ALLOCATION_NOT_FOUND"
	CodeAllocationInvalid  = "ERR_ALLOCATION_INVALID"
	CodeAllocationMismatch = "ERR_ALLOCATION_MISMATCH"
	CodeInvalidProportion  = "ERR_INVALID_PROPORTION"

	CodeFileNotFound      = "ERR_FILE_NOT_FOUND"
	CodeFileInvalid       = "ERR_FILE_INVALID"
	CodeFileTooLarge      = "ERR_FILE_TOO_LARGE"
	CodeUnsupportedFormat = "ERR_UNSUPPORTED_FORMAT"

	CodeDatabaseError     = "ERR_DATABASE"
	CodeTransactionFailed = "ERR_TRANSACTION_FAILED"
	CodeDuplicateEntry    = "ERR_DUPLICATE_ENTRY"
)

type AppError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	HTTPStatus int                    `json:"-"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Err        error                  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
func (e *AppError) Unwrap() error {
	return e.Err
}
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}
func New(code, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}
func Wrap(err error, code, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}
func ErrInternal(message string) *AppError {
	return New(CodeInternal, message, http.StatusInternalServerError)
}
func ErrValidation(message string) *AppError {
	return New(CodeValidation, message, http.StatusBadRequest)
}
func ErrNotFound(message string) *AppError {
	return New(CodeNotFound, message, http.StatusNotFound)
}
func ErrConflict(message string) *AppError {
	return New(CodeConflict, message, http.StatusConflict)
}
func ErrBadRequest(message string) *AppError {
	return New(CodeBadRequest, message, http.StatusBadRequest)
}
func ErrUnauthorized(message string) *AppError {
	return New(CodeUnauthorized, message, http.StatusUnauthorized)
}
func ErrForbidden(message string) *AppError {
	return New(CodeForbidden, message, http.StatusForbidden)
}
func ErrTooManyRequests(message string) *AppError {
	return New(CodeTooManyRequests, message, http.StatusTooManyRequests)
}
func ErrInvalidCredentials() *AppError {
	return New(CodeInvalidCredentials, "Invalid email or password", http.StatusUnauthorized)
}
func ErrTokenExpired() *AppError {
	return New(CodeTokenExpired, "Token has expired", http.StatusUnauthorized)
}
func ErrTokenInvalid() *AppError {
	return New(CodeTokenInvalid, "Invalid token", http.StatusUnauthorized)
}
func ErrOTPInvalid() *AppError {
	return New(CodeOTPInvalid, "Invalid OTP code", http.StatusBadRequest)
}
func ErrOTPExpired() *AppError {
	return New(CodeOTPExpired, "OTP code has expired", http.StatusBadRequest)
}
func ErrPINInvalid() *AppError {
	return New(CodePINInvalid, "Invalid PIN", http.StatusBadRequest)
}
func ErrPINLocked() *AppError {
	return New(CodePINLocked, "PIN is locked due to too many failed attempts", http.StatusForbidden)
}
func ErrPINRequired() *AppError {
	return New(CodePINRequired, "PIN verification required for this operation", http.StatusForbidden)
}
func ErrUserNotFound() *AppError {
	return New(CodeUserNotFound, "User not found", http.StatusNotFound)
}
func ErrUserAlreadyExists() *AppError {
	return New(CodeUserAlreadyExists, "User with this email already exists", http.StatusConflict)
}
func ErrUserNotApproved() *AppError {
	return New(CodeUserNotApproved, "User registration is pending approval", http.StatusForbidden)
}
func ErrUserSuspended() *AppError {
	return New(CodeUserSuspended, "User account is suspended", http.StatusForbidden)
}
func ErrProfileIncomplete() *AppError {
	return New(CodeProfileIncomplete, "User profile is incomplete", http.StatusForbidden)
}
func ErrTenantNotFound() *AppError {
	return New(CodeTenantNotFound, "Tenant not found", http.StatusNotFound)
}
func ErrTenantInactive() *AppError {
	return New(CodeTenantInactive, "Tenant is inactive", http.StatusForbidden)
}
func ErrPermissionDenied() *AppError {
	return New(CodePermissionDenied, "You do not have permission to perform this action", http.StatusForbidden)
}
func ErrRoleNotFound() *AppError {
	return New(CodeRoleNotFound, "Role not found", http.StatusNotFound)
}
func ErrAccessForbidden(message string) *AppError {
	return New(CodeAccessForbidden, message, http.StatusForbidden)
}
func ErrPlatformAdminRequired() *AppError {
	return New(CodePlatformAdminRequired, "This operation requires platform administrator privileges", http.StatusForbidden)
}
func ErrEmployeeNotFound() *AppError {
	return New(CodeEmployeeNotFound, "Employee not found", http.StatusNotFound)
}
func ErrEmployeeExists() *AppError {
	return New(CodeEmployeeExists, "Employee with this NIK already exists", http.StatusConflict)
}
func ErrContributionNotFound() *AppError {
	return New(CodeContributionNotFound, "Contribution not found", http.StatusNotFound)
}
func ErrAllocationNotFound() *AppError {
	return New(CodeAllocationNotFound, "Allocation not found", http.StatusNotFound)
}
func ErrInvalidProportion() *AppError {
	return New(CodeInvalidProportion, "Allocation proportions must sum to 100%", http.StatusBadRequest)
}
func ErrFileNotFound() *AppError {
	return New(CodeFileNotFound, "File not found", http.StatusNotFound)
}
func ErrFileTooLarge(maxSize string) *AppError {
	return New(CodeFileTooLarge, fmt.Sprintf("File exceeds maximum size of %s", maxSize), http.StatusBadRequest)
}
func ErrUnsupportedFormat(format string) *AppError {
	return New(CodeUnsupportedFormat, fmt.Sprintf("Unsupported file format: %s", format), http.StatusBadRequest)
}
func ErrDatabase(message string) *AppError {
	return New(CodeDatabaseError, message, http.StatusInternalServerError)
}
func ErrDuplicateEntry(field string) *AppError {
	return New(CodeDuplicateEntry, fmt.Sprintf("Duplicate entry for %s", field), http.StatusConflict)
}
func Is(err, target error) bool {
	return errors.Is(err, target)
}
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}
func GetHTTPStatus(err error) int {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}
func GetCode(err error) string {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.Code
	}
	return CodeInternal
}
