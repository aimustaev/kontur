package api

import (
	"github.com/aimustaev/service-workflow/internal/manager_workflow"
)

// ConfigVersionRepository определяет интерфейс для работы с конфигурациями
type ConfigVersionRepository interface {
	// GetLatestActive возвращает последнюю активную версию конфигурации по имени
	GetLatestActive(name string) (*manager_workflow.ConfigVersion, error)

	// GetByVersion возвращает конкретную версию конфигурации
	GetByVersion(name, version string) (*manager_workflow.ConfigVersion, error)

	// Create создает новую версию конфигурации
	Create(config *manager_workflow.ConfigVersion) error

	// Update обновляет существующую версию конфигурации
	Update(config *manager_workflow.ConfigVersion) error

	// List возвращает список версий конфигураций по фильтру
	List(filter manager_workflow.ConfigVersionFilter) ([]*manager_workflow.ConfigVersion, error)

	// Deactivate деактивирует версию конфигурации
	Deactivate(name, version string) error
}
