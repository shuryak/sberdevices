package oauth

import (
	"log"

	"github.com/shuryak/sberdevices/internal/config"
	"github.com/shuryak/sberdevices/internal/oauth"
)

// Handlers - https://datatracker.ietf.org/doc/html/rfc6749#section-4.1
type Handlers struct {
	flow *oauth.CodeFlowWithOTP
	cfg  *config.Config
	log  *log.Logger
}

func NewHandlers(flow *oauth.CodeFlowWithOTP, cfg *config.Config, log *log.Logger) *Handlers {
	return &Handlers{
		flow: flow,
		cfg:  cfg,
		log:  log,
	}
}
