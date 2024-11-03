package oauth

import (
	"context"
	"net/http"

	"github.com/shuryak/sberdevices/internal/api"
)

type startReq struct {
	Phone string `query:"phone"`
}

func (req startReq) Validate(_ *api.Context) error {
	// TODO: just uncomment
	// if len(req.Phone) != 11 {
	// 	return errors.New("invalid phone number")
	// }

	return nil
}

type StartResp struct {
	Code string `json:"code"`
}

func (h *Handlers) Start(_ *api.Context, req *startReq) (*StartResp, int) {
	code, err := h.flow.Start(context.Background(), req.Phone)
	if err != nil {
		h.log.Printf("failed to start flow, err: %v\n", err)
		return nil, http.StatusInternalServerError
	}

	return &StartResp{Code: code}, http.StatusOK
}
