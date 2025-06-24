package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/joomcode/errorx"

	errs "crud_api/internal/errors"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	code := http.StatusInternalServerError
	typ := errorx.GetTypeName(err)
	msg := err.Error()

	switch {
	case errorx.IsOfType(err, errs.ErrUserNotFound):
		code = http.StatusNotFound
	case errorx.IsOfType(err, errs.ErrUserDuplicate):
		code = http.StatusConflict
	case errorx.IsOfType(err, errs.ErrInvalidPayload):
		code = http.StatusBadRequest
	case errorx.IsOfType(err, errs.ErrUnauthorized):
		code = http.StatusUnauthorized
	case errorx.IsOfType(err, errs.ErrDbFailure):
		code = http.StatusInternalServerError
		msg = "Internal server error"
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Code: code, Type: typ, Message: msg})
	
}
