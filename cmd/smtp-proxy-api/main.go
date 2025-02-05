package main

import (
	"context"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sinland/smtp-proxy-api/internal/config"
	http_handler "github.com/sinland/smtp-proxy-api/internal/handler/http"
	http_middleware "github.com/sinland/smtp-proxy-api/internal/middleware/http"
	"github.com/sinland/smtp-proxy-api/internal/service/smtp"
	"github.com/sinland/smtp-proxy-api/internal/service/tg"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	ctx := context.Background()

	configPath := flag.String("config", "config.yml", "path to config file")
	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		slog.ErrorContext(ctx, "failed to load file", "error", err)
		return
	}

	smtpService := newSmtpService(cfg)
	tgService, err := tg.NewService(tg.Config{BotToken: cfg.Server.BotToken})
	if err != nil {
		slog.ErrorContext(ctx, "failed to init tg service", "error", err)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	r := chi.NewRouter()
	r.Use(http_middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	server := http_handler.NewMainController(http_handler.MainControllerConfig{AppConfig: cfg, MailSender: smtpService, TgService: tgService})

	go func() {
		if err = server.Run(r); err != nil {
			slog.ErrorContext(ctx, "failed to run main controller", "error", err)
		}
	}()
	slog.Info("server started", "port", cfg.Server.Port)

	<-stop
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err = server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error: %v", err)
	}
}

func newSmtpService(cfg *config.Config) *smtp.Service {
	return smtp.NewService(cfg.SMTP.Server, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password)
}
