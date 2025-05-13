package activity

import (
	"context"

	"github.com/aimustaev/service-workflow/internal/config"
)

var globalConfig *config.Config

// SetConfig sets the global configuration
func SetConfig(cfg *config.Config) {
	globalConfig = cfg
}

// GetConfig returns the global configuration
func GetConfig() *config.Config {
	return globalConfig
}

// ContextWithConfig creates a new context with config
func ContextWithConfig(cfg *config.Config) context.Context {
	return context.WithValue(context.Background(), "config", cfg)
}

// GetConfigFromContext gets config from context
func GetConfigFromContext(ctx context.Context) *config.Config {
	if cfg, ok := ctx.Value("config").(*config.Config); ok {
		return cfg
	}
	return nil
}
