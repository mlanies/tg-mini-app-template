package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// SendWebAppMessageHandler - обработчик для отправки сообщений из веб-приложения
func SendWebAppMessageHandler(bot *gotgbot.Bot, appURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос на отправку сообщения из веб-приложения")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Функция отправки сообщения из веб-приложения"))
	}
}

// UsersHandler - обработчик для работы с пользователями
func UsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос на получение списка пользователей")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Функция обработки пользователей"))
	}
}

// AppointmentsHandler - обработчик для работы с записями
func AppointmentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос на получение записей")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Функция обработки записей"))
	}
}

// ServicesHandler - обработчик для работы с услугами
func ServicesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос на получение услуг")

		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		// Пример получения списка услуг из базы данных
		rows, err := db.Query("SELECT id, name, price, duration FROM services")
		if err != nil {
			log.Printf("Ошибка при выполнении запроса к базе данных: %v", err)
			http.Error(w, "Ошибка при выполнении запроса к базе данных", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var services []Service
		for rows.Next() {
			var service Service
			if err := rows.Scan(&service.ID, &service.Name, &service.Price, &service.Duration); err != nil {
				log.Printf("Ошибка при сканировании строки: %v", err)
				http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
				return
			}
			services = append(services, service)
		}

		// Проверка на наличие ошибок после цикла rows.Next()
		if err = rows.Err(); err != nil {
			log.Printf("Ошибка при чтении строк: %v", err)
			http.Error(w, "Ошибка при чтении данных", http.StatusInternalServerError)
			return
		}

		// Установка заголовков и возврат списка услуг в формате JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(services); err != nil {
			log.Printf("Ошибка при кодировании ответа в JSON: %v", err)
			http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		}
	}
}

// Структура услуги
type Service struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Duration int     `json:"duration"` // Продолжительность услуги в минутах
}
