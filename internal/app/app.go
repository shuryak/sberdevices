package app

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/shuryak/sberdevices/internal/adapter"
	"github.com/shuryak/sberdevices/internal/api"
	"github.com/shuryak/sberdevices/internal/api/devices"
	oauthApi "github.com/shuryak/sberdevices/internal/api/oauth"
	"github.com/shuryak/sberdevices/internal/config"
	"github.com/shuryak/sberdevices/internal/oauth"
	postgresPkg "github.com/shuryak/sberdevices/internal/pkg/postgres"
	"github.com/shuryak/sberdevices/internal/pkg/router"
	"github.com/shuryak/sberdevices/internal/pkg/smarthome/auth"
	"github.com/shuryak/sberdevices/internal/pkg/smarthome/client"
	"github.com/shuryak/sberdevices/internal/storage/postgres"
)

func Run() {
	logger := log.Default()

	cfg, err := config.Read("./config.yaml")
	if err != nil {
		logger.Fatalf("read config, err: %v\n", err)
		return
	}

	slogLogger := initLogger()

	pg, err := postgresPkg.New(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database,
		),
		postgresPkg.NewOptions().SetMaxPoolSize(3),
		slogLogger,
	)
	if err != nil {
		// TODO: handle
		panic(err)
	}

	smartHomeClient := client.NewClient(20*time.Second, logger)

	codeStorage := postgres.NewOAuthCodeStorage(pg, cfg.Auth.CodeLength)
	sessionStorage := postgres.NewSessionStorage(pg)

	authorizer := auth.NewAuthorizer(logger)
	flow := oauth.NewCodeFlowWithOTP(
		codeStorage, sessionStorage, postgres.NewUnitOfWorkFactory(pg, sessionStorage, codeStorage),
		adapter.NewAuthorizer(authorizer), cfg.Auth.AccessTokenLength, cfg.Auth.RefreshTokenLength,
	)

	r := router.New(logger)

	oauthHandlers := oauthApi.NewHandlers(flow, cfg, logger)
	oauthGroup := router.NewGroup("/oauth",
		router.GET("/start", oauthHandlers.Authorize),
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

func initLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
