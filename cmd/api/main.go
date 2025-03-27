package main

import (
	"context"
	"log/slog"

	"github.com/zetsub0/yakvs/internal/adapters/tarantool"
	"github.com/zetsub0/yakvs/internal/app/api"
	"github.com/zetsub0/yakvs/internal/app/http"
	"github.com/zetsub0/yakvs/internal/config"
	"github.com/zetsub0/yakvs/internal/logger"
	"github.com/zetsub0/yakvs/internal/modules/manager"
)

func main() {

	ctx := context.Background()

	cfg := config.ParseConfig()

	store := tarantool.New(ctx, cfg.Tarantool)

	log := logger.SetupLogger(cfg.Env)
	slog.SetDefault(log)

	mng := manager.New(store)

	mux := api.NewHandler(api.NewAPI(mng))

	srv := http.New(cfg.HttpServer, mux)

	srv.Run(ctx)

}
