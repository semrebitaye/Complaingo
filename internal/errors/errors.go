package errors

import "github.com/joomcode/errorx"

var (
	commonErrors = errorx.NewNamespace("common")

	ErrUserNotFound   = errorx.NewType(commonErrors, "user_not_found")
	ErrUserDuplicate  = errorx.NewType(commonErrors, "user_duplicate")
	ErrInvalidPayload = errorx.NewType(commonErrors, "invalid_payload")
	ErrUnauthorized   = errorx.NewType(commonErrors, "unauthorized")
	ErrDbFailure      = errorx.NewType(commonErrors, "db_failure")
)
