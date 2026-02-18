package main

import (
	"database/sql"
	"log"
	"time"

	"go-gin-ticketing-backend/internal/config"
	"go-gin-ticketing-backend/seed"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("mysql", cfg.LocalhostDatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	gofakeit.Seed(time.Now().UnixNano())

	if err := seed.Run(db); err != nil {
		log.Fatal(err)
	}
}
