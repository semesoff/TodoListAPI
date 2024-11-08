package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"todo-list-api/config"
)

var (
	db *sql.DB
)

func InitDB() {
	var err error
	cfg := config.GetConfig().Database
	db, err = sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
		return
	}
	if db.Ping() != nil {
		log.Fatalf("Error pinging database: %q", err)
		return
	}
	log.Println("Database is initialized")
}

func GetDB() *sql.DB {
	return db
}
