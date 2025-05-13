package activity

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

// ProcessMessageActivity обрабатывает входящее сообщение
func (a *Activity) ProcessMessageActivity(ctx context.Context, message string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Обработка сообщения", "message", message)
	return nil
}

// WaitActivity выполняет ожидание указанное количество секунд
func (a *Activity) WaitActivity(ctx context.Context, seconds int) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Начало ожидания", "seconds", seconds)
	time.Sleep(time.Duration(seconds) * time.Second)
	logger.Info("Ожидание завершено")
	return nil
}
