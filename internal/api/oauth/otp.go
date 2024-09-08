package oauth

import (
	"context"
	"errors"
	"net/http"

	"github.com/shuryak/sberhack/internal/api"
)

type otpReq struct {
	Code string `query:"code"`
	OTP  string `query:"otp"`
}

func (req otpReq) Validate(_ *api.Context) error {
	if len(req.OTP) != 5 {
		return errors.New("invalid otp")
	}

	return nil
}

type OTPResp struct {
	OK bool `json:"ok"`
}

func (h *Handlers) OTP(_ *api.Context, req *otpReq) (*OTPResp, int) {
	_, err := h.flow.CreateSession(context.Background(), req.Code, req.OTP)
	if err != nil {
		return &OTPResp{OK: false}, http.StatusBadRequest
	}

	return &OTPResp{OK: true}, http.StatusOK
}
