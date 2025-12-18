package main

import (
	"log"

	"ticket-io/internal/config"
	"ticket-io/internal/database"

	userhandler "ticket-io/internal/user/handler"

	statusservice "ticket-io/internal/user/service/status"
	userservice "ticket-io/internal/user/service/user"

	statusrepository "ticket-io/internal/user/repository/status"
	userrepository "ticket-io/internal/user/repository/user"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()

	db, err := database.NewMysql(cfg.DockerDatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	api := r.Group("/api/v1")

	// repositories
	userRepo := userrepository.New(db)
	statusRepo := statusrepository.New(db)

	// services
	statusService := statusservice.New(statusRepo)
	userService := userservice.New(userRepo, statusService)

	// handlers & routes
	userhandler.RegisterRoutes(api, userService)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
	})

	r.Run(":8080")
}
