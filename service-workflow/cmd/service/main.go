package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.temporal.io/sdk/client"

	"github.com/aimustaev/service-workflow/internal/api"
	"github.com/aimustaev/service-workflow/internal/config"
	"github.com/aimustaev/service-workflow/internal/manager_workflow"
	"github.com/aimustaev/service-workflow/internal/usecase"
)

func main() {
	cfg := config.Load()

	// Подключаемся к PostgreSQL
	db, err := sqlx.Connect("postgres", cfg.GetPostgresDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализируем репозиторий конфигураций
	configRepo := manager_workflow.NewPostgresConfigRepository(db)

	// Создаем клиент Temporal
	c, err := client.NewClient(client.Options{
		HostPort: cfg.GetTemporalAddr(),
	})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	// Инициализируем use cases
	startWorkflowUseCase := usecase.NewStartWorkflowUseCase(c)
	startWorkflowV2UseCase := usecase.NewStartV2WorkflowUseCase(c)

	// Создаем HTTP хендлеры
	startHandler := api.NewStartWorkflowHandler(startWorkflowUseCase)
	startV2Handler := api.NewStartV2WorkflowHandler(startWorkflowV2UseCase)
	healthHandler := &api.HealthHandler{}

	// Создаем хендлеры для конфигураций
	getLatestConfigHandler := api.NewGetLatestConfigHandler(configRepo)
	getVersionConfigHandler := api.NewGetVersionConfigHandler(configRepo)
	createConfigHandler := api.NewCreateConfigHandler(configRepo)
	updateConfigHandler := api.NewUpdateConfigHandler(configRepo)
	listConfigHandler := api.NewListConfigHandler(configRepo)
	deactivateConfigHandler := api.NewDeactivateConfigHandler(configRepo)

	// Создаем роутер
	router := mux.NewRouter()

	// Регистрируем маршруты для воркфлоу
	router.HandleFunc("/start", startHandler.Handle).Methods(http.MethodPost)
	router.HandleFunc("/startV2", startV2Handler.Handle).Methods(http.MethodPost)
	router.HandleFunc("/health", healthHandler.Handle).Methods(http.MethodGet)

	// Регистрируем маршруты для конфигураций
	router.HandleFunc("/config/{name}/latest", getLatestConfigHandler.Handle).Methods(http.MethodGet)
	router.HandleFunc("/config/{name}/version/{version}", getVersionConfigHandler.Handle).Methods(http.MethodGet)
	router.HandleFunc("/config", createConfigHandler.Handle).Methods(http.MethodPost)
	router.HandleFunc("/config/{name}/version/{version}", updateConfigHandler.Handle).Methods(http.MethodPut)
	router.HandleFunc("/config/{name}", listConfigHandler.Handle).Methods(http.MethodGet)
	router.HandleFunc("/config/{name}/version/{version}/deactivate", deactivateConfigHandler.Handle).Methods(http.MethodPost)

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:    cfg.GetHTTPAddr(),
		Handler: router,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Starting HTTP server on %s", cfg.GetHTTPAddr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Ждем сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Даем 5 секунд на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
