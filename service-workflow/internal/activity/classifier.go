package activity

import (
	"context"

	"go.temporal.io/sdk/activity"
)

// ClassifierAcitivity классифицирует
func ClassifierAcitivity(ctx context.Context, message string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Классификация тикета")
	return nil
}
