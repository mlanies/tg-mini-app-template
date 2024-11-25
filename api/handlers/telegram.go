package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// InitBot - инициализация бота Telegram
func InitBot(botToken string) *gotgbot.Bot {
	bot, err := gotgbot.NewBot(botToken, nil)
	if err != nil {
		log.Printf("Ошибка инициализации Telegram Bot API: %v", err)
		return nil
	}
	log.Println("Telegram Bot API успешно инициализирован")
	return bot
}

// CreateBotEndpointHandler - обработчик для взаимодействия с Telegram Bot API
func CreateBotEndpointHandler(bot *gotgbot.Bot, appURL string, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Обработка запроса %s", r.URL.Path)
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		var update gotgbot.Update
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			log.Printf("Ошибка декодирования обновления: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if update.Message == nil {
			log.Printf("Обновление не содержит сообщения")
			http.Error(w, "Обновление бота не содержит сообщения", http.StatusBadRequest)
			return
		}

		userID := update.Message.From.Id
		username := update.Message.From.Username
		firstName := update.Message.From.FirstName
		lastName := update.Message.From.LastName
		languageCode := update.Message.From.LanguageCode

		_, err = db.Exec(`INSERT INTO users (id, username, first_name, last_name, language_code, created_at) VALUES ($1, $2, $3, $4, $5, NOW()) ON CONFLICT (id) DO NOTHING`, userID, username, firstName, lastName, languageCode)
		if err != nil {
			log.Printf("Ошибка записи пользователя в базу данных: %v", err)
		} else {
			log.Printf("Данные пользователя %s успешно сохранены", username)
		}

		message := "💅 Добро пожаловать в наш салон красоты, " + firstName
		if lastName != "" {
			message += " " + lastName
		}
		message += "!\n\n✨ Забронируйте услугу маникюра или педикюра прямо сейчас, нажав на кнопку ниже. Воспользуйтесь нашим мини-приложением, чтобы увидеть весь спектр услуг и удобное расписание."

		opts := &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{{Text: "📅 Записаться", WebApp: &gotgbot.WebAppInfo{Url: appURL}}}},
			},
		}

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
