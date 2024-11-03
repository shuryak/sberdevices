package oauth

import (
	"github.com/shuryak/sberdevices/internal/api"
	"github.com/shuryak/sberdevices/internal/api/common"
)

func (h *Handlers) ErrHandler(_ *api.Context, err error) (interface{}, int) {
	return common.NewErrorResp(err)
}
