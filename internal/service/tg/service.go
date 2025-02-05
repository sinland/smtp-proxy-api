package tg

import (
	"context"
	"github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"log/slog"
)

type Service struct {
	bot *bot.Bot
}

func NewService(cfg Config) (*Service, error) {
	svc := &Service{}

	opts := []bot.Option{
		bot.WithDefaultHandler(svc.defaultHandler),
	}

	b, err := bot.New(cfg.BotToken, opts...)
	if err != nil {
		return nil, err
	}

	svc.bot = b

	return svc, nil
}

func (s *Service) Start(ctx context.Context) {
	s.bot.Start(ctx)
}

type Config struct {
	BotToken string
}

func (s *Service) SendMessage(ctx context.Context, to, message string) error {
	msg, err := s.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    to,
		Text:      message,
		ParseMode: tg_models.ParseModeHTML,
	})
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "sent telegram message", "msg_id", msg.ID)

	return nil
}

func (s *Service) defaultHandler(ctx context.Context, b *bot.Bot, update *tg_models.Update) {
	slog.InfoContext(ctx, "handling telegram update",
		"update_id", update.ID,
	)
}
