package main

import (
	"log"
	"net/http"
	"os"

	"tg-mini-app-template/api/database"
	"tg-mini-app-template/api/handlers"
)

func main() {
	log.Println("Запуск сервиса API")

	// Инициализация базы данных
	db := database.InitDB()
	defer db.Close()

	// Получение переменных окружения
	webAppURL := os.Getenv("TELEGRAM_WEB_APP_URL")
	if webAppURL == "" {
		log.Fatal("TELEGRAM_WEB_APP_URL не задан")
	}
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не задан")
	}

	// Инициализация Telegram Bot
	bot := handlers.InitBot(botToken)
	if bot == nil {
		log.Fatalf("Ошибка инициализации Telegram Bot API")
	}

	// Настройка маршрутов
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/api/bot", handlers.CreateBotEndpointHandler(bot, webAppURL, db))
	mux.HandleFunc("/api/send-webapp", handlers.SendWebAppMessageHandler(bot, webAppURL))
	mux.HandleFunc("/api/users", handlers.UsersHandler(db))
	mux.HandleFunc("/api/appointments", handlers.AppointmentsHandler(db))
	mux.HandleFunc("/api/services", handlers.ServicesHandler(db))
	mux.HandleFunc("/api/masters", handlers.MastersHandler(db))

	// Оберните маршруты в CORS Middleware
	handler := handlers.CorsMiddleware(mux)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
