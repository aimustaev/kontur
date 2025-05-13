package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	httphandler "github.com/aimustaev/service-gateway/internal/http"
	rpcClient "github.com/aimustaev/service-gateway/internal/tickets"
)

func main() {
	// Загружаем переменные окружения из .env файла
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	if err := godotenv.Load(".env." + env); err != nil {
		log.Printf("Warning: .env.%s file not found: %v", env, err)
	}

	// Создаем клиент для tickets сервиса
	ticketsClient, err := rpcClient.NewClient()
	if err != nil {
		log.Fatalf("Failed to create tickets client: %v", err)
	}

	// Создаем HTTP хендлер
	handler := httphandler.NewHandler(ticketsClient)

	// Настраиваем роутер
	router := mux.NewRouter()
	router.HandleFunc("/api/tickets", handler.GetAllTickets).Methods("GET")
	router.HandleFunc("/api/ticket/{id}/messages", handler.GetTicketMessages).Methods("GET")

	// Получаем порт для HTTP сервера из переменных окружения или используем значение по умолчанию
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// Запускаем HTTP сервер
	log.Printf("Starting HTTP server on :%s", httpPort)
	if err := http.ListenAndServe(":"+httpPort, router); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
