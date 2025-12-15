package main

import (
	"database/sql"
	"log"
	"time"

	"learning-gin/internal/config"
	"learning-gin/internal/seed"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("mysql", cfg.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	gofakeit.Seed(time.Now().UnixNano())

	if err := seed.SeedUsersTable(db, 100); err != nil {
		log.Fatal(err)
	}
}
