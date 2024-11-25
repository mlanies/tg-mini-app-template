package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func InitBot(botToken string) *gotgbot.Bot {
	bot, err := gotgbot.NewBot(botToken, nil)
	if err != nil {
		log.Printf("Ошибка инициализации Telegram Bot API: %v", err)
		return nil
	}
	log.Println("Telegram Bot API успешно инициализирован")
	return bot
}

func CreateBotEndpointHandler(bot *gotgbot.Bot, appURL string, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Обработка запроса %s", r.URL.Path)
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		// Декодирование обновления от Telegram
		var update gotgbot.Update
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			log.Printf("Ошибка декодирования обновления: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Проверка на наличие сообщения в обновлении
		if update.Message == nil {
			log.Printf("Обновление не содержит сообщения")
			http.Error(w, "Обновление бота не содержит сообщения", http.StatusBadRequest)
			return
		}

		// Получение информации о пользователе
		userID := update.Message.From.Id
		username := update.Message.From.Username
		firstName := update.Message.From.FirstName
		lastName := update.Message.From.LastName
		languageCode := update.Message.From.LanguageCode

		// Запись информации о пользователе в базу данных
		_, err = db.Exec(`INSERT INTO users (id, username, first_name, last_name, language_code, created_at) VALUES ($1, $2, $3, $4, $5, NOW()) ON CONFLICT (id) DO NOTHING`,
			userID, username, firstName, lastName, languageCode)
		if err != nil {
			log.Printf("Ошибка записи пользователя в базу данных: %v", err)
		} else {
			log.Printf("Данные пользователя %s успешно сохранены", username)
		}

		// Формирование приветственного сообщения
		message := "💅 Добро пожаловать в наш салон красоты, " + firstName
		if lastName != "" {
			message += " " + lastName
		}
		message += "!\n\n✨ Забронируйте услугу маникюра или педикюра прямо сейчас, нажав на кнопку ниже. Воспользуйтесь нашим мини-приложением, чтобы увидеть весь спектр услуг и удобное расписание."

		// Формирование кнопки для перехода в веб-приложение
		opts := &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{Text: "📅 Записаться", WebApp: &gotgbot.WebAppInfo{Url: appURL}},
					},
				},
			},
		}

		// Отправка сообщения пользователю
		resp, err := bot.SendMessage(update.Message.Chat.Id, message, opts)
		if err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
			http.Error(w, "Не удалось отправить сообщение: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Сообщение успешно отправлено: %+v", resp)
		w.WriteHeader(http.StatusOK)
	}
}

// Обработчик для получения мастеров
func MastersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос на получение мастеров")

		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		// Получение списка мастеров из базы данных
		rows, err := db.Query("SELECT id, name, experience FROM masters")
		if err != nil {
			log.Printf("Ошибка при выполнении запроса к базе данных: %v", err)
			http.Error(w, "Ошибка при выполнении запроса к базе данных", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Создание списка мастеров
		var masters []Master
		for rows.Next() {
			var master Master
			if err := rows.Scan(&master.ID, &master.Name, &master.Experience); err != nil {
				log.Printf("Ошибка при сканировании строки: %v", err)
				http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
				return
			}
			masters = append(masters, master)
		}

		// Проверка на наличие ошибок после цикла rows.Next()
		if err = rows.Err(); err != nil {
			log.Printf("Ошибка при чтении строк: %v", err)
			http.Error(w, "Ошибка при чтении данных", http.StatusInternalServerError)
			return
		}

		// Установка заголовков и возврат списка мастеров в формате JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(masters); err != nil {
			log.Printf("Ошибка при кодировании ответа в JSON: %v", err)
			http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		}
	}
}

// Структура мастера
type Master struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Experience int    `json:"experience"` // Поле опыта работы мастера
}
