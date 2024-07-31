package main

import (
	"github.com/q2rd/gRPC_sso_go/internal/app"
	"github.com/q2rd/gRPC_sso_go/internal/config"
	"github.com/q2rd/gRPC_sso_go/internal/custom_logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	logger := custom_logger.SetupLogger(cfg.Env)
	logger.Info(
		"Start app with:",
		slog.String("env", cfg.Env),
		slog.Int("port:", cfg.GRPC.Port),
	)
	logger.Debug("check logger performance")
	application := app.NewApp(logger, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GRPCSrv.Run()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	signalFromChan := <-stop
	logger.Info("Application stopping: ", slog.String("signal", signalFromChan.String()))
	application.GRPCSrv.Stop()
	logger.Info("Application stopped.")
}
