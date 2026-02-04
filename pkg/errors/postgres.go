package errors

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func TranslatePostgres(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return SentinelNotFound
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return SentinelDuplicate
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return SentinelDuplicate
		case "23503":
			return SentinelForeignKey
		case "40001", "40P01":
			return SentinelConflict
		}
	}

	if isPostgresConnectionError(err) {
		return err
	}

	return err
}

func isPostgresConnectionError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	patterns := []string{
		"connection refused",
		"connection reset",
		"broken pipe",
		"no connection",
	}
	for _, p := range patterns {
		if strings.Contains(msg, p) {
			return true
		}
	}
	return false
}
