package app

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/shuryak/sberdevices/internal/adapter"
	"github.com/shuryak/sberdevices/internal/api"
	"github.com/shuryak/sberdevices/internal/api/devices"
	oauthApi "github.com/shuryak/sberdevices/internal/api/oauth"
	"github.com/shuryak/sberdevices/internal/config"
	"github.com/shuryak/sberdevices/internal/oauth"
	"github.com/shuryak/sberdevices/internal/storage/mem"
	"github.com/shuryak/sberdevices/pkg/router"
	"github.com/shuryak/sberdevices/pkg/smarthome/auth"
	"github.com/shuryak/sberdevices/pkg/smarthome/client"
)

func Run() {
	logger := log.Default()

	cfg := config.Read(logger, "./config.yaml")

	r := router.New(logger)

	smartHomeClient := client.NewClient(20*time.Second, logger)

	codeStorage := mem.NewOAuthCodeStorage(64)
	sessionStorage := mem.NewSessionStorage(128, 128)

	authorizer := auth.NewAuthorizer(logger)
	flow := oauth.NewCodeFlowWithOTP(codeStorage, sessionStorage, adapter.NewAuthorizer(authorizer))

	oauthHandlers := oauthApi.NewHandlers(flow, logger)
	oauthGroup := router.NewGroup("/oauth",
		router.GET("/start", oauthHandlers.Start),
		router.GET("/otp", oauthHandlers.OTP),
		router.POST("/token", oauthHandlers.Token),
		router.POST("/refresh", oauthHandlers.Refresh),
	).SetPreHandler(api.CORS).SetErrHandler(oauthHandlers.ErrHandler)

	devicesHandlers := devices.NewHandlers(smartHomeClient, flow)
	devicesGroup := router.NewGroup("/user/devices",
		router.GET("", devicesHandlers.Devices),
		router.POST("/query", devicesHandlers.DevicesQuery),
		router.POST("/action", devicesHandlers.DevicesAction),
	).SetPreHandler(devicesHandlers.Tokens)

	r.Add(router.NewGroup("/api", oauthGroup, devicesGroup))
	r.Add(router.NewGroup("/api/v1.0", oauthGroup, devicesGroup))

	port, _ := strings.CutPrefix(cfg.Server.Port, ":")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("listening on %s\n", srv.Addr)

	startServer(srv)
}

func startServer(srv *http.Server) {
	err := srv.ListenAndServe()
	if err != nil {
		panic(err) // TODO: handle
	}
}
