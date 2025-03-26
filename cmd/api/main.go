package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/zetsub0/yakvs/internal/adapters/tarantool"
	"github.com/zetsub0/yakvs/internal/app/api"
	"github.com/zetsub0/yakvs/internal/app/http"
	"github.com/zetsub0/yakvs/internal/config"
	"github.com/zetsub0/yakvs/internal/modules/manager"
)

func main() {

	ctx := context.Background()

	cfg := config.ParseConfig()

	store := tarantool.New(ctx, cfg.Tarantool)

	logger := setupLogger(cfg.Env)

	slog.SetDefault(logger)

	mng := manager.New(store)

	mux := api.NewHandler(api.NewAPI(mng))

	srv := http.New(cfg.HttpServer, mux)

	srv.Run(ctx)

}

func setupLogger(env string) *slog.Logger {

	opts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "file" {
				return slog.Attr{}
			}
			return a
		},
	}

	switch env {
	case "dev":
		opts.Level = slog.LevelDebug
	case "local":
		opts.Level = slog.LevelDebug
	case "prod":
		opts.Level = slog.LevelError
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	logger := slog.New(handler)
	return logger
}
