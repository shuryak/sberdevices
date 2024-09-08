package common

import (
	"errors"
	"net/http"
)

type ErrorResp struct {
	Error struct {
		Code string `json:"code"`
	} `json:"error"`
}

func NewErrorResp(err error) (*ErrorResp, int) {
	statusCode := http.StatusBadRequest
	if v, ok := errToHTTPStatusCode[err]; ok {
		statusCode = v
	}

	return &ErrorResp{
		Error: struct {
			Code string `json:"code"`
		}{
			Code: err.Error(),
		},
	}, statusCode
}

var (
	ErrInvalidRequest = errors.New("invalid_request")
)

var errToHTTPStatusCode = map[error]int{
	ErrInvalidRequest: http.StatusBadRequest,
}
