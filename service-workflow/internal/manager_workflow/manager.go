package manager_workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aimustaev/service-workflow/internal/engine"
)

// ConfigManager manages workflow configurations with dynamic updates
type ConfigManager struct {
	repo           ConfigVersionRepository
	cache          map[string]*cachedConfig
	cacheMutex     sync.RWMutex
	updateInterval time.Duration
	stopChan       chan struct{}
}

type cachedConfig struct {
	config     *ConfigVersion
	definition engine.WorkflowDefinition
	schema     *json.RawMessage
	updatedAt  time.Time
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(repo ConfigVersionRepository, updateInterval time.Duration) *ConfigManager {
	return &ConfigManager{
		repo:           repo,
		cache:          make(map[string]*cachedConfig),
		updateInterval: updateInterval,
		stopChan:       make(chan struct{}),
	}
}

// Start starts the configuration manager
func (m *ConfigManager) Start(ctx context.Context) {
	log.Printf("Starting ConfigManager with update interval: %v", m.updateInterval)
	go m.updateLoop(ctx)
}

// Stop stops the configuration manager
func (m *ConfigManager) Stop() {
	close(m.stopChan)
}

// GetWorkflowDefinition returns the latest workflow definition for the given name
func (m *ConfigManager) GetWorkflowDefinition(name string) (engine.WorkflowDefinition, error) {
	log.Printf("Getting workflow definition for: %s", name)
	m.cacheMutex.RLock()
	cached, exists := m.cache[name]
	m.cacheMutex.RUnlock()

	if !exists {
		log.Printf("Cache miss for workflow: %s, loading from database", name)
		return m.loadAndCacheConfig(name)
	}

	// Проверяем, не устарела ли конфигурация
	if time.Since(cached.updatedAt) > m.updateInterval {
		log.Printf("Config for %s is stale (age: %v), updating in background", name, time.Since(cached.updatedAt))
		// Пытаемся обновить конфигурацию в фоне
		go func() {
			if err := m.updateConfig(name); err != nil {
				log.Printf("Failed to update config for %s: %v", name, err)
			}
		}()
	} else {
		log.Printf("Using cached config for %s (age: %v)", name, time.Since(cached.updatedAt))
	}

	return cached.definition, nil
}

// GetWorkflowSchema возвращает визуальную схему для воркфлоу
func (m *ConfigManager) GetWorkflowSchema(name string) (*json.RawMessage, error) {
	m.cacheMutex.RLock()
	cached, exists := m.cache[name]
	m.cacheMutex.RUnlock()

	if !exists {
		config, err := m.repo.GetLatestActive(name)
		if err != nil {
			return nil, err
		}
		if config == nil {
			return nil, ErrConfigNotFound
		}

		def, err := m.parseConfig(config)
		if err != nil {
			return nil, err
		}

		m.cacheMutex.Lock()
		m.cache[name] = &cachedConfig{
			config:     config,
			definition: def,
			schema:     config.Schema,
			updatedAt:  time.Now(),
		}
		m.cacheMutex.Unlock()

		return config.Schema, nil
	}

	return cached.schema, nil
}

// loadAndCacheConfig загружает конфигурацию из БД и кэширует её
func (m *ConfigManager) loadAndCacheConfig(name string) (engine.WorkflowDefinition, error) {
	config, err := m.repo.GetLatestActive(name)
	if err != nil {
		return engine.WorkflowDefinition{}, err
	}
	if config == nil {
		return engine.WorkflowDefinition{}, ErrConfigNotFound
	}

	def, err := m.parseConfig(config)
	if err != nil {
		return engine.WorkflowDefinition{}, err
	}

	m.cacheMutex.Lock()
	m.cache[name] = &cachedConfig{
		config:     config,
		definition: def,
		updatedAt:  time.Now(),
	}
	m.cacheMutex.Unlock()

	return def, nil
}

// updateConfig обновляет конфигурацию в кэше
func (m *ConfigManager) updateConfig(name string) error {
	config, err := m.repo.GetLatestActive(name)
	if err != nil {
		return err
	}
	if config == nil {
		return ErrConfigNotFound
	}

	m.cacheMutex.RLock()
	cached, exists := m.cache[name]
	m.cacheMutex.RUnlock()

	// Если конфигурация не изменилась, пропускаем обновление
	if exists && cached.config.ID == config.ID {
		return nil
	}

	def, err := m.parseConfig(config)
	if err != nil {
		return err
	}

	m.cacheMutex.Lock()
	m.cache[name] = &cachedConfig{
		config:     config,
		definition: def,
		updatedAt:  time.Now(),
	}
	m.cacheMutex.Unlock()

	log.Printf("Updated workflow configuration for %s to version %s", name, config.Version)
	return nil
}

// updateLoop периодически проверяет обновления конфигураций
func (m *ConfigManager) updateLoop(ctx context.Context) {
	log.Printf("Starting update loop with interval: %v", m.updateInterval)
	ticker := time.NewTicker(m.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("ConfigManager update loop stopped: context cancelled")
			return
		case <-m.stopChan:
			log.Printf("ConfigManager update loop stopped: stop signal received")
			return
		case <-ticker.C:
			log.Printf("Running periodic config update check")
			m.cacheMutex.RLock()
			names := make([]string, 0, len(m.cache))
			for name := range m.cache {
				names = append(names, name)
			}
			m.cacheMutex.RUnlock()

			log.Printf("Checking updates for workflows: %v", names)
			for _, name := range names {
				if err := m.updateConfig(name); err != nil {
					log.Printf("Failed to update config for %s: %v", name, err)
				}
			}
		}
	}
}

// parseConfig парсит конфигурацию в определение воркфлоу
func (m *ConfigManager) parseConfig(config *ConfigVersion) (engine.WorkflowDefinition, error) {
	var def engine.WorkflowDefinition
	if err := json.Unmarshal(config.Content, &def); err != nil {
		return engine.WorkflowDefinition{}, fmt.Errorf("failed to unmarshal workflow definition: %w", err)
	}
	return def, nil
}

var (
	ErrConfigNotFound = errors.New("configuration not found")
)
