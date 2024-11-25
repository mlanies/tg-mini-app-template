package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Структура пользователя
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Структура записи
type Appointment struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	ServiceID       int    `json:"service_id"`
	AppointmentTime string `json:"appointment_time"`
}

// Структура услуги
type Service struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Duration int     `json:"duration"` // Продолжительность услуги в минутах
}

// Структура мастера
type Master struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Experience int    `json:"experience"`
}

// HealthHandler - обработчик для проверки состояния сервера
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// CorsMiddleware - добавляет заголовки CORS к ответу
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// UsersHandler - обработчик для получения пользователей
func UsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос на получение списка пользователей")
		// Логика для обработки запроса пользователей
	}
}

// AppointmentsHandler - обработчик для работы с записями
func AppointmentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос на получение записей")
		// Логика для обработки запроса на записи
	}
}

// ServicesHandler - обработчик для работы с услугами
func ServicesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос на получение услуг")
		// Логика для обработки запроса услуг
	}
}

// MastersHandler - обработчик для работы с мастерами
func MastersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос на получение мастеров")
		// Логика для обработки запроса мастеров
	}
}

// ServicesByMasterHandler - обработчик для получения услуг мастера по его ID
func ServicesByMasterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос на получение услуг мастера")

		masterIDStr := r.URL.Query().Get("master_id")
		if masterIDStr == "" {
			http.Error(w, "Не указан master_id", http.StatusBadRequest)
			return
		}

		masterID, err := strconv.Atoi(masterIDStr)
		if err != nil {
			http.Error(w, "Некорректный master_id", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT id, name, price, duration FROM services WHERE master_id = $1", masterID)
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

		if err = rows.Err(); err != nil {
			log.Printf("Ошибка при чтении строк: %v", err)
			http.Error(w, "Ошибка при чтении данных", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(services); err != nil {
			log.Printf("Ошибка при кодировании ответа в JSON: %v", err)
			http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		}
	}
}
