package oauth

import (
	"log"

	"github.com/shuryak/sberdevices/internal/oauth"
)

type Handlers struct {
	flow *oauth.CodeFlowWithOTP
	log  *log.Logger
}

func NewHandlers(flow *oauth.CodeFlowWithOTP, log *log.Logger) *Handlers {
	return &Handlers{
		flow: flow,
		log:  log,
	}
}
