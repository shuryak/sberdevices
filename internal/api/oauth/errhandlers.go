package oauth

import (
	"github.com/shuryak/sberhack/internal/api"
	"github.com/shuryak/sberhack/internal/api/common"
)

func (h *Handlers) ErrHandler(_ *api.Context, err error) (interface{}, int) {
	return common.NewErrorResp(err)
}
