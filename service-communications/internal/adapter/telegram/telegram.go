package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"

	"github.com/aimustaev/service-communications/internal/adapter"
)

type Problem struct {
	children []Problem
	id       int64
	label    string
}

// TelegramAdapter implements the Adapter interface for Telegram
type TelegramAdapter struct {
	botToken     string
	chatID       string
	logger       *logrus.Logger
	db           adapter.Database
	bot          *tgbotapi.BotAPI
	lastUpdateID int64
}

// NewTelegramAdapter creates a new Telegram adapter
func NewTelegramAdapter(botToken, chatID string, logger *logrus.Logger, db adapter.Database) *TelegramAdapter {
	return &TelegramAdapter{
		botToken:     botToken,
		chatID:       chatID,
		logger:       logger,
		db:           db,
		lastUpdateID: 0,
	}
}

// Connect establishes a connection to Telegram
func (a *TelegramAdapter) Connect(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(a.botToken)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	a.bot = bot
	a.logger.Infof("Connected to Telegram bot: %s", bot.Self.UserName)
	return nil
}

// Disconnect closes the connection to Telegram
func (a *TelegramAdapter) Disconnect(ctx context.Context) error {
	// No connection to close for HTTP-based service
	return nil
}

// GetMessages retrieves messages from Telegram
func (a *TelegramAdapter) GetMessages(ctx context.Context) ([]adapter.Message, error) {

	if a.bot == nil {
		return nil, fmt.Errorf("bot not initialized")
	}

	updateConfig := tgbotapi.NewUpdate(int(a.lastUpdateID))
	updates, err := a.bot.GetUpdates(updateConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %w", err)
	}

	var messages []adapter.Message
	for _, update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		if update.Message != nil {
			message := adapter.Message{
				ID:      fmt.Sprintf("%d", update.Message.MessageID),
				From:    update.Message.From.UserName,
				To:      a.bot.Self.UserName,
				Subject: "Telegram Message",
				Body:    update.Message.Text,
				Tags:    []string{"telegram"},
				Channel: "telegram",
			}

			messages = append(messages, message)

			if err := a.db.SaveEmail(ctx, message); err != nil {
				a.logger.Errorf("Failed to save email %s to database: %v", message.ID, err)
			}
		}

		a.lastUpdateID = int64(update.UpdateID) + 1
	}

	return messages, nil
}

func (a *TelegramAdapter) SendButtons(problems []Problem) error {
	var buttons []tgbotapi.InlineKeyboardButton
	if len(problems) == 0 {
		msg := tgbotapi.NewMessage(a.getChatID(), "Это конец")
		a.bot.Send(msg)
		return nil
	}

	for _, problem := range problems {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(problem.label, fmt.Sprintf("%d", problem.id)))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
	msg := tgbotapi.NewMessage(a.getChatID(), "Выберите действие:")
	msg.ReplyMarkup = keyboard
	_, err := a.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// MarkAsProcessed marks a message as processed
func (a *TelegramAdapter) MarkAsProcessed(ctx context.Context, messageID string) error {
	a.logger.Infof("Marked message %s as processed", messageID)
	return nil
}

// getChatID converts the chat ID string to int64
func (a *TelegramAdapter) getChatID() int64 {
	var chatID int64
	fmt.Sscanf(a.chatID, "%d", &chatID)
	return chatID
}
