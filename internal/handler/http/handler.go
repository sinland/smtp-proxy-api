package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sinland/smtp-proxy-api/internal/config"
	http_middleware "github.com/sinland/smtp-proxy-api/internal/middleware/http"
	"github.com/sinland/smtp-proxy-api/internal/service/smtp"
	"github.com/sinland/smtp-proxy-api/internal/service/tg"
	"net/http"
)

type MainController struct {
	appConfig  *config.Config
	server     *http.Server
	mailSender *smtp.Service
	tgService  *tg.Service
}

func (c *MainController) Run(mux *chi.Mux) error {
	c.initialize(mux)

	c.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", c.appConfig.Server.Port),
		Handler: mux,
	}
	if err := c.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (c *MainController) Shutdown(ctx context.Context) error {
	return c.server.Shutdown(ctx)
}

func (c *MainController) initialize(mux *chi.Mux) {
	mux.Group(func(r chi.Router) {
		r.Use(middleware.Throttle(10))
		r.Use(http_middleware.NewApiKeyMiddleware(c.appConfig.Server.ApiKey))
		r.Post("/auth/token", c.Login)
	})
	mux.Group(func(r chi.Router) {
		r.Use(http_middleware.NewJWTMiddleware(c.appConfig.Server.JwtSecret))

		// user state endpoints
		r.Post("/email/send-message", c.SendEmailMessage)
		r.Post("/tg/send-message", c.SendTelegramMessage)
	})
}

type MainControllerConfig struct {
	AppConfig  *config.Config
	MailSender *smtp.Service
	TgService  *tg.Service
}

func NewMainController(cfg MainControllerConfig) *MainController {
	return &MainController{
		appConfig:  cfg.AppConfig,
		mailSender: cfg.MailSender,
		tgService:  cfg.TgService,
	}
}
