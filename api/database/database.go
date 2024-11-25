package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() *sql.DB {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка проверки соединения с базой данных: %v", err)
	}

	log.Println("Соединение с базой данных успешно установлено")
	return db
}
